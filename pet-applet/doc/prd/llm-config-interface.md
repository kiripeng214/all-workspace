---
status: draft
author: loop-task
created: 2026-06-20
---

# PRD: LLM 配置移到 config + 抽接口

## 1. 背景与动机

- LLM 配置（API Key / URL / Model）目前硬编码用 os.Getenv，与项目配置系统不一致
- 模型调用逻辑混在 llm.go 中 OpenAI/Anthropic 两种格式 if-else 拼凑，难以扩展
- 新加 provider（如 Ollama 本地模型）需要改已有代码

## 2. 需求描述

### 2.1 功能列表

| # | 功能 | 优先级 | 说明 |
|---|------|--------|------|
| 1 | config.yaml 增加 llm 配置段 | P0 | api_key / api_url / model / provider |
| 2 | LLMProvider 接口 | P0 | `Ask(ctx, prompt) (string, error)` |
| 3 | OpenAI 实现 | P0 | 现有逻辑抽为独立类型 |
| 4 | 通过 config 初始化 provider | P0 | main → config → knowledge → provider |

## 3. 验收标准

- [ ] config.yaml 有 llm 配置段（含默认值注释）
- [ ] config-local.yaml 可覆盖 llm 配置
- [ ] LLMProvider 接口定义清晰
- [ ] `QueryLLM` 通过 provider 调用，不直接读 env
- [ ] `go vet ./...` 通过
- [ ] `go test ./...` 通过

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v0.1 | 初稿 | loop-task |
