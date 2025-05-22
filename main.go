// main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"codeprobot/internal/config"
	"codeprobot/internal/watcher"
	"github.com/fsnotify/fsnotify"
)

func main() {
	cfg, err := config.Load(".github/agent.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("🤖 CodePRobot is running...")
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
