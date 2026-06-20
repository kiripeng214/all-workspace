---
status: draft
author: loop-task
created: 2026-06-20
based-on: doc/prd/llm-config-interface.md
---

# Plan: LLM 配置移到 config + 抽接口

## 1. 方案概述

- config 层：新增 LLMConfig 结构体 + yaml 映射
- 接口层：`LLMProvider` 接口，`Ask(ctx, prompt) (string, error)`
- 实现层：`OpenAIProvider` + `AnthropicProvider`
- 知识库层：`QueryLLM` 通过 provider 调用，config 从 Init 传入

## 2. 涉及文件

| 文件 | 改动类型 | 改动内容 |
|------|---------|---------|
| `backend/config/config.go` | 修改 | 新增 LLMConfig 结构体 |
| `backend/config/config.yaml` | 修改 | 新增 llm 配置段 |
| `backend/config/config-local.yaml` | 修改 | 新增 llm 配置段注释 |
| `backend/knowledge/knowledge.go` | 修改 | Init 接受 LLMConfig，创建 provider |
| `backend/knowledge/llm.go` | 重构 | 抽接口 + OpenAI 实现 + 移除 os.Getenv |
| `backend/knowledge/anthropic.go` | 新增 | Anthropic provider 实现 |
| `backend/main.go` | 修改 | 传递 cfg.LLM 给 knowledge.Init |

## 3. 实施步骤

### 步骤 1: 扩展 config
- Config 增加 `LLM LLMConfig`
- LLMConfig: Provider, APIKey, APIURL, Model

### 步骤 2: 重构 llm.go
- 定义 `LLMProvider` 接口
- 提取 `OpenAIProvider` 类型，实现接口
- `QueryLLM` 接收 provider 参数

### 步骤 3: 新建 anthropic.go
- `AnthropicProvider` 实现接口

### 步骤 4: 更新调用链
- knowledge.Init 接收 LLMConfig
- main.go 传递 cfg.LLM

### 步骤 5: 验证
- `go vet ./...` + `go test ./...`

## 4. 风险与回退

| 风险 | 概率 | 应对 |
|------|------|------|
| API Key 在 config.yaml 中提交 | 低 | config-local.yaml 覆盖，config.yaml 留空 |

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v0.1 | 初稿 | loop-task |
