# 🐾 宠物助手 — 简历项目描述

## 项目名称

**宠物助手** — 全栈宠物管理 + AI 知识库应用

## 一句话概括

基于 Go + Vue 3 的全栈宠物管理平台，集成 RAG 知识库实现宠物养护智能问答。配置 GitHub CI、Docker 容器化，支持多种 LLM 提供商切换。

## 技术栈

Go / Gin / MySQL / Goose / Docker  
Vue 3 / TypeScript / uni-app / Vite 5  
chromem-go / Ollama / LLM API / GitHub Actions

## 项目职责

**后端开发**（Go + Gin + MySQL）
- 设计并实现 15 个 RESTful API，支持宠物、喂养计划、喂养记录 CRUD
- 使用 Goose 管理数据库版本迁移，支持 up/down 回滚
- 分层配置系统：默认值 → config.yaml → 本地覆盖 config-local.yaml
- 多阶段 Docker 构建，最终镜像约 15MB

**RAG 知识库**（向量检索 + LLM）
- 基于 chromem-go 构建内存向量库，实现语义相似度搜索
- 混合搜索策略：向量检索 + 关键词精确过滤，无关结果不返回
- 设计 LLMProvider 接口，支持 OpenAI / DeepSeek / Anthropic / Ollama 切换
- 集成 Ollama 本地 embedding（nomic-embed-text），API 不可用时自动降级
- 22 条专业知识覆盖狗/猫/赤狐的饲养、健康、行为

**前端开发**（Vue 3 + TypeScript + uni-app）
- 组件化重构，将 449 行大组件拆分为 3 个子组件，提升可维护性
- 23 个前端测试用例保障核心交互质量

**工程化**
- GitHub Actions CI：3 个并行 job（go vet+test、vue-tsc+vitest、docker build）
- 56 个后端测试 + 性能基准测试（embedding 240ns/op）
- conventional commit 规范，PRD→Plan→Code 标准化交付流程

## 项目地址

https://github.com/kiripeng214/all-workspace
