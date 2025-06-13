package main

import (
	"flag"
	"fmt"
	"os"
)
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

	// Show help message if requested or if required flags are missing
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

