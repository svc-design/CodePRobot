APP_NAME := codepilot
MAIN_FILE := main.go
TAG ?= latest
CONFIG ?= .github/agent.yaml

.PHONY: all build run clean init watch pr help

all: build

init:
	GOPROXY=https://goproxy.cn,direct go get gopkg.in/yaml.v2
	GOPROXY=https://goproxy.cn,direct go get github.com/fsnotify/fsnotify
	GOPROXY=https://goproxy.cn,direct go get github.com/spf13/cobra@latest
	go mod tidy

build:
	go build -o $(APP_NAME) $(MAIN_FILE)

run:
	go run $(MAIN_FILE)

watch:
	go run $(MAIN_FILE) watch

pr:
	go run $(MAIN_FILE) pr --tag=$(TAG)

clean:
	rm -f $(APP_NAME)

help:
	@echo "🤖 CodePRobot CLI Usage"
	@echo ""
	@echo "make build             编译 CodePRobot 二进制"
	@echo "make run               启动主逻辑（监听 + PR 自动化）"
	@echo "make watch             启动文件监听模式"
	@echo "make pr TAG=xxx        模拟手动触发一个指定 TAG 的 PR"
	@echo "make init              初始化依赖（go mod tidy + cobra）"
	@echo "make clean             删除构建产物"
