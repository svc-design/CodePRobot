package cmd

import (
	"log"

	"codeprobot/internal/config"
	"codeprobot/internal/generator"
	"codeprobot/internal/github"
	"codeprobot/internal/gitops"
	"codeprobot/internal/watcher"
)

func main() {
	// 加载配置
	cfg, err := config.Load("example/agent.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化模块
	w := watcher.NewWatcher(cfg.WatchPaths)
	github := watcher.NewGitHubClient(cfg.GitHub.Token, cfg.GitHub.Repo)
	gitOps := watcher.NewGitOps(cfg.GitHub.Repo)
	generator := watcher.NewGenerator(cfg.OpenAI.APIKey, cfg.OpenAI.Model, cfg.OpenAI.Temperature)

	// 开始监听
	log.Println("Starting CodePRobot Agent...")
	for {
		events := w.Check() // 检查文件变动
		for _, event := range events {
			log.Printf("Detected event: %v", event)

			// 触发关键字检查
			if event.ContainsKeywords(cfg.TriggerKeywords) {
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
				if err := github.CreatePullRequest(branchName, cfg.GitHub.BaseBranch, cfg.GitHub.Reviewer); err != nil {
					log.Printf("Error creating PR: %v", err)
					continue
				}
			}
		}
	}
}
