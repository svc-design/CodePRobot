// main.go
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"codeprobot/internal/config"
	"codeprobot/internal/watcher"
	"github.com/fsnotify/fsnotify"
)

func main() {
	// 解析命令行参数
	verbose := flag.Bool("v", false, "enable verbose mode")
	help := flag.Bool("help", false, "show help message")
	once := flag.Bool("once", false, "run once in demo mode")
	flag.Parse()

	if *help {
		log.Println("Usage: ./codepilot [--once] [-v] [--help]")
		log.Println("  --once   run once in demo mode (non-daemon)")
		log.Println("  -v       enable verbose logging")
		log.Println("  --help   show this help message")
		return
	}

	if *verbose {
		log.Println("🔍 Verbose mode enabled.")
	}

	cfg, err := config.Load("example/agent.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if *once {
		log.Println("🚀 Running CodePRobot in once/demo mode...")
		log.Println("✅ Simulate file change on README.md")
		log.Println("🧠 Triggering ChatGPT generator...")
		log.Println("🔧 GitOps: commit and push")
		log.Println("📬 GitHub: create PR")
		log.Println("✅ Demo completed")
		return
	}

	log.Println("🤖 CodePRobot is running in watch mode...")
	w := watcher.NewWatcher(cfg.WatchPaths)
	defer w.Close()

	go w.Start(func(event fsnotify.Event) {
		log.Printf("🧠 Triggering pipeline for %s", event.Name)
		// TODO: 触发 generator, gitops, github 等模块处理
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("👋 CodePRobot stopped.")
}
