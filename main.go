package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Repository struct {
	Name string `json:"name"`
	FullName string `json:"full_name"`
}

type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
}

type CommitStats struct {
	Repository string
	Count      int
	Commits    []Commit
}

type GitHubClient struct {
	Token    string
	Username string
	Client   *http.Client
}

func NewGitHubClient(token, username string) *GitHubClient {
	return &GitHubClient{
		Token:    token,
		Username: username,
		Client:   &http.Client{Timeout: 30 * time.Second},
	}
}

// makeRequest makes an authenticated HTTP request to the GitHub API
func (g *GitHubClient) makeRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+g.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func(g *GitHubClient) GetUserRepositories() ([]Repository, error){
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&sort=updated", g.Username)
	body, err := g.makeRequest(url)
	if err != nil {
		return nil, err
	}
	var repos []Repository
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, err
	}
return repos, nil
}

func (g *GitHubClient) GetCommitsForRepo(repo Repository, since time.Time) ([]Commit, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/commits?author=%s&since=%s&per_page=100",
		repo.FullName, g.Username, since.Format(time.RFC3339))
	body, err := g.makeRequest(url)
	if err != nil {
		return nil, err
	}

	var commits []Commit
	if err := json.Unmarshal(body, &commits); err != nil {
		return nil, err
	}
	return commits, nil
}

func getTimeRanges() (time.Time, time.Time, time.Time, time.Time) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday is 0, but we treat it as 7 for our calculation
		weekday = 7
	}
	weekStart := todayStart.AddDate(0, 0, -(weekday - 1))
	weekEnd := weekStart.Add(7 * 24 * time.Hour)

	return todayStart, todayEnd, weekStart, weekEnd
}

func filterCommitsByTimeRange(commits []Commit, start, end time.Time) []Commit {
	var filtered []Commit
	for _, commit := range commits {
		commitTime := commit.Commit.Author.Date
		if !commitTime.Before(start) && commitTime.Before(end) {
			filtered = append(filtered, commit)
		}
	}
	return filtered
}

// printCommitStats displays commit statistics in a formatted way
func printCommitStats(title string, stats []CommitStats, totalCommits int) {
	fmt.Printf("\n=== %s ===\n", title)
	fmt.Printf("Total commits: %d\n", totalCommits)
	if totalCommits == 0 {
		fmt.Println("No commits found for this period.")
		return
	}
	fmt.Println("\nBy repository:")

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Count > stats[j].Count
	})

	for _, stat := range stats {
		if stat.Count > 0 {
			fmt.Printf("  %s: %d commits\n", stat.Repository, stat.Count)
			maxShow := 3
			if len(stat.Commits) < maxShow {
				maxShow = len(stat.Commits)
			}
			for i := 0; i < maxShow; i++ {
				commit := stat.Commits[i]
				message := strings.Split(commit.Commit.Message, "\n")[0] // First line only
				if len(message) > 60 {
					message = message[:57] + "..."
				}
				fmt.Printf("    - %s\n", message)
			}
			if len(stat.Commits) > maxShow {
				fmt.Printf("    ... and %d more commits\n", len(stat.Commits)-maxShow)
			}
		}
	}
}

func main(){
	var(
		token = flag.String("token", "", "Gitub Personal Access Token")
		username = flag.String("username", "", "Github username")
		help = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()
	if *token == "" {
		*token = os.Getenv("GITHUB_TOKEN")
	}
	if *username == "" {
		*username = os.Getenv("GITHUB_USERNAME")
	}

	if *help || *token == "" || *username == "" {
		fmt.Println("GitHub Commit Tracker")
		fmt.Println("Usage: go run main.go -token <your_token> -username <your_username>")
		fmt.Println("Alternatively, set GITHUB_TOKEN and GITHUB_USERNAME environment variables.")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
		return
	}

	client := NewGitHubClient(*token, *username)
	fmt.Printf("Fetching commit statistics for %s...\n", *username)

	repos, err := client.GetUserRepositories()
	if err != nil {
		fmt.Printf("Error fetching repositories: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d repositories to check.\n", len(repos))

	todayStart, todayEnd, weekStart, weekEnd := getTimeRanges()
	var todayStats, weekStats []CommitStats
	var todayTotal, weekTotal int = 0, 0

	for i, repo := range repos {
		fmt.Printf("\rProcessing repository %d/%d: %s      ", i+1, len(repos), repo.Name)
		commits, err := client.GetCommitsForRepo(repo, weekStart)
		if err != nil {
			continue
		}

		todayCommits := filterCommitsByTimeRange(commits, todayStart, todayEnd)
		if len(todayCommits) > 0 {
			todayStats = append(todayStats, CommitStats{
				Repository: repo.Name,
				Count:      len(todayCommits),
				Commits:    todayCommits,
			})
			todayTotal += len(todayCommits)
		}

		weekCommits := commits
		if len(weekCommits) > 0 {
			weekStats = append(weekStats, CommitStats{
				Repository: repo.Name,
				Count:      len(weekCommits),
				Commits:    weekCommits,
			})
			weekTotal += len(weekCommits)
		}
}
fmt.Println("\n\nProcessing complete.") 
	printCommitStats("TODAY'S COMMITS", todayStats, todayTotal)
	printCommitStats("THIS WEEK'S COMMITS", weekStats, weekTotal)

	fmt.Printf("\nTime period (Today): %s\n", todayStart.Format("Jan 2, 2006"))
	fmt.Printf("Time period (This Week): %s - %s\n", weekStart.Format("Jan 2"), weekEnd.Add(-time.Second).Format("Jan 2, 2006"))
}

