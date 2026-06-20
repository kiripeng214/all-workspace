---
status: draft
author: loop-task
created: 2026-06-20
based-on: doc/prd/pet-knowledge-rag.md
---

# Plan: 宠物知识库 (RAG + chromem-go + LLM)

## 1. 方案概述

- 并行模式：后端（go-dev）+ 前端（uni-app-dev）同时开发
- 后端：chromem-go 向量库 + seed 知识数据 + LLM API 调用
- 前端：知识页面 + 详情页入口

## 2. 涉及文件

| 文件 | 改动类型 | 改动内容 |
|------|---------|---------|
| `backend/knowledge/knowledge.go` | 新增 | chromem-go 初始化 + 向量检索 |
| `backend/knowledge/seed.go` | 新增 | 宠物知识种子数据 |
| `backend/knowledge/router.go` | 新增 | /api/knowledge/search 路由 |
| `backend/main.go` | 修改 | 注册知识库路由 |
| `miniprogram/src/api/knowledge.ts` | 新增 | 知识库 API 客户端 |
| `miniprogram/src/pages/knowledge/index.vue` | 新增 | 知识库展示页面 |
| `miniprogram/src/pages/pets/detail.vue` | 修改 | 增加"知识"入口按钮 |
| `pages.json` | 修改 | 注册知识库路由 |

## 3. 实施步骤

### 并行组 A
- [go-dev] Step 1: knowledge 包（chromem-go + 种子数据）
- [uni-app-dev] Step 2: 前端知识页面 + API 层

### 并行组 B（依赖 Step 1）
- [go-dev] Step 3: 注册路由 + main.go

## 4. 风险与回退

| 风险 | 概率 | 应对 |
|------|------|------|
| 无 LLM API Key | 中 | 环境变量 `LLM_API_KEY` 控制，缺失时返回纯检索结果 |
| chromem-go 兼容性 | 低 | 纯 Go 实现，无外部依赖 |
| embedding 需要 API | 中 | 可用 chromem-go 内置的 OpenAI embedding 或 HuggingFace |

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v0.1 | 初稿 | loop-task |
