
# GitPeek
## GitHub Commit Tracker CLI

A command-line tool written in **Go** that fetches and displays:
- The number of commits you’ve made today and this week
- Commits grouped by repository

## Features

✅ Fetches your repositories using GitHub API  
✅ Lists your commits filtered by today/this week  
✅ Shows commit messages and count  

---

## Requirements

- [Go](https://go.dev/dl/) 1.18 or newer  
- A **GitHub personal access token**

---


## Usage

### 1️⃣ Set up your GitHub token and username in your terminal:

```bash
# Replace YOUR_TOKEN_HERE and YOUR_USERNAME_HERE
export GITHUB_TOKEN=YOUR_TOKEN_HERE
export GITHUB_USERNAME=YOUR_USERNAME_HERE
```

👉 To make this permanent, add those lines to your `~/.bashrc` or `~/.bash_profile`:

```bash
echo 'export GITHUB_TOKEN=YOUR_TOKEN_HERE' >> ~/.bashrc
echo 'export GITHUB_USERNAME=YOUR_USERNAME_HERE' >> ~/.bashrc
source ~/.bashrc
```

If you use Zsh, update `~/.zshrc` instead.

---

### 2️⃣ Run the app

```bash
go run main.go
```

or build it:

```bash
go build -o gitpeek main.go
./gitpeek
```

---

## Example Output

```
=== TODAY'S COMMITS ===
Total commits: 3

By repository:
  my-repo: 3 commits
    - Fix bug in authentication flow
    - Add unit tests for login
    - Refactor token handler

=== THIS WEEK'S COMMITS ===
Total commits: 12
...
```

---

## Notes

- The tool uses the GitHub API v3. Make sure your token has at least `repo` scope for private repos.
- It fetches up to 100 commits per repo by default (GitHub API limit per request).
- TUI support can be extended using [Bubble Tea](https://github.com/charmbracelet/bubbletea).

---

## Contributing

Pull requests and feature suggestions are welcome!
