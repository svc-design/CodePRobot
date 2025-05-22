package watcher

import (
	"fmt"
	"os/exec"
)

type GitHubClient struct {
	token string
	repo  string
}

func NewGitHubClient(token, repo string) *GitHubClient {
	return &GitHubClient{token: token, repo: repo}
}

func (g *GitHubClient) CreatePullRequest(branch, baseBranch string, reviewers []string) error {
	title := fmt.Sprintf("Auto-generated PR for branch %s", branch)
	cmd := exec.Command("gh", "pr", "create", "--base", baseBranch, "--head", branch, "--title", title)
	return cmd.Run()
}
