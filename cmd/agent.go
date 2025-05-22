package main

import (
	"log"
	"os"

    "codeprobot/internal/config"
    "codeprobot/internal/watcher"
    "codeprobot/internal/github"
    "codeprobot/internal/gitops"
    "codeprobot/internal/generator"
)

func main() {
	// 加载配置
	config, err := internal.LoadConfig(".github/agent.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化模块
	watcher := internal.NewWatcher(config.WatchPaths)
	github := internal.NewGitHubClient(config.GitHub.Token, config.GitHub.Repo)
	gitOps := internal.NewGitOps(config.GitHub.Repo)
	generator := internal.NewGenerator(config.OpenAI.APIKey, config.OpenAI.Model, config.OpenAI.Temperature)

	// 开始监听
	log.Println("Starting CodePilot Agent...")
	for {
		events := watcher.Check() // 检查文件变动
		for _, event := range events {
			log.Printf("Detected event: %v", event)

			// 触发关键字检查
			if event.ContainsKeywords(config.TriggerKeywords) {
				log.Println("Trigger keyword detected. Processing...")

				// 调用生成器
				code, err := generator.Generate(event.Content)
				if err != nil {
					log.Printf("Error generating code: %v", err)
					continue
				}

				// Git 操作
				branchName := "auto/" + event.Name
				if err := gitOps.CreateBranch(branchName); err != nil {
					log.Printf("Error creating branch: %v", err)
					continue
				}
				if err := gitOps.CommitAndPush(branchName, code); err != nil {
					log.Printf("Error in git commit or push: %v", err)
					continue
				}

				// 创建 PR
				if err := github.CreatePullRequest(branchName, config.GitHub.BaseBranch, config.GitHub.Reviewer); err != nil {
					log.Printf("Error creating PR: %v", err)
					continue
				}
			}
		}
	}
}
