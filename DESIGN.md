# CodePRobot 设计草案

本文档描述了扩展 CodePRobot 的几个关键功能，以满足如下需求：

1. **支持接入 GitHub 仓库**：利用 `gh` CLI 与 GitHub API 交互，自动创建分支、提交并发起 PR。
2. **支持使用 `codex-cli` 与 `claude-code`**：在生成代码时除了调用 OpenAI API，还可以通过执行 `npx codex-cli` 或 `npx claude-code` 获得结果。
3. **支持配置代理访问大模型**：在配置文件中新增 `proxy.http` 字段，程序启动时若设置则配置 HTTP 代理。
4. **从 GitHub Issue 读取需求并循环提交**：新增 `IssueReader` 组件，通过 `gh issue view` 拉取 Issue 内容，作为生成器输入。根据 `loop_count` 重复执行生成、提交和 PR 流程直到构建成功。

上述功能均已在 `example/agent.yaml` 中展示配置方式，具体实现可在 `internal/watcher` 目录下查看。
