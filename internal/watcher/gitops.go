package watcher

import (
	"os"
	"os/exec"
)

type GitOps struct {
	repo string
}

func NewGitOps(repo string) *GitOps {
	return &GitOps{repo: repo}
}

func (g *GitOps) CreateBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	return cmd.Run()
}

func (g *GitOps) CommitAndPush(branchName, content string) error {
	// 写入文件
	err := os.WriteFile("generated_code.go", []byte(content), 0644)
	if err != nil {
		return err
	}

	// Git 提交
	cmd := exec.Command("git", "add", ".")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git", "commit", "-m", "Generated code by CodePilot")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Git 推送
	cmd = exec.Command("git", "push", "-u", "origin", branchName)
	return cmd.Run()
}
