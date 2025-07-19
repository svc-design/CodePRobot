package watcher

import (
	"os/exec"
)

type IssueReader struct {
	repo string
}

func NewIssueReader(repo string) *IssueReader {
	return &IssueReader{repo: repo}
}

func (i *IssueReader) Read(issueNumber string) (string, error) {
	cmd := exec.Command("gh", "issue", "view", issueNumber, "--repo", i.repo, "--json", "body", "-q", ".body")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
