package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"github.com/fsnotify/fsnotify"
	"codeprobot/internal/config"
	"codeprobot/internal/generator"
	"codeprobot/internal/github"
	"codeprobot/internal/gitops"
	"codeprobot/internal/watcher"
)

func main() {
	cfg, err := config.Load(".github/agent.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	w := watcher.NewWatcher(cfg.WatchPaths)
	gh := watcher.NewGitHubClient(cfg.GitHub.Token, cfg.GitHub.Repo)
	gitOps := watcher.NewGitOps(cfg.GitHub.Repo)
	gen := watcher.NewGenerator(os.Getenv("OPENAI_API_KEY"), cfg.OpenAI.Model, cfg.OpenAI.Temperature)

	log.Println("Starting CodePilot Agent...")
	w.Start(func(event fsnotify.Event) {
		log.Printf("Detected event: %v", event)

		if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
			return
		}

		data, err := os.ReadFile(event.Name)
		if err != nil {
			log.Printf("Error reading file: %v", err)
			return
		}

		if !containsKeywords(string(data), cfg.TriggerKeywords) {
			return
		}

		log.Println("Trigger keyword detected. Processing...")

		code, err := gen.Generate(string(data))
		if err != nil {
			log.Printf("Error generating code: %v", err)
			return
		}

		branchName := "auto/" + filepath.Base(event.Name)
		if err := gitOps.CreateBranch(branchName); err != nil {
			log.Printf("Error creating branch: %v", err)
			return
		}
		if err := gitOps.CommitAndPush(branchName, code); err != nil {
			log.Printf("Error in git commit or push: %v", err)
			return
		}

		if err := gh.CreatePullRequest(branchName, "main", cfg.GitHub.Reviewer); err != nil {
			log.Printf("Error creating PR: %v", err)
		}
	})
}

func containsKeywords(content string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(content, kw) {
			return true
		}
	}
	return false
}
