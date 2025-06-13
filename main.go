package main

import (
	"encoding/json"
	"io"
	"net/http"
	"flag"
	"fmt"
	"time"
	"os"
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
	fmt.Printf("%s", *token)
	fmt.Printf("%s", *username)
}

