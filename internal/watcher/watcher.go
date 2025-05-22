package watcher

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

// Callback 是当文件变动发生时调用的函数类型
type Callback func(event fsnotify.Event)

// Watcher 封装了对文件系统的监听功能
type Watcher struct {
	paths   []string
	watcher *fsnotify.Watcher
}

// NewWatcher 创建一个新的 Watcher 实例
func NewWatcher(paths []string) *Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("❌ Error creating watcher: %v", err)
	}

	for _, path := range paths {
		err := watcher.Add(path)
		if err != nil {
			log.Printf("⚠️ Warning: could not watch path %s: %v", path, err)
		} else {
			log.Printf("✅ Watching path: %s", path)
		}
	}

	return &Watcher{
		paths:   paths,
		watcher: watcher,
	}
}

// Start 开始监听，并在事件发生时调用回调函数
func (w *Watcher) Start(callback Callback) {
	log.Println("👀 Watcher started...")
	for {
		select {
		case event := <-w.watcher.Events:
			log.Printf("📄 File event: %s %s", event.Op, event.Name)
			callback(event)

		case err := <-w.watcher.Errors:
			log.Printf("❌ Watcher error: %v", err)
		}
	}
}

// Close 关闭 Watcher
func (w *Watcher) Close() {
	if w.watcher != nil {
		_ = w.watcher.Close()
	}
}
