# 🐾 宠物助手 — 技能展示

> Go + Vue 3 + TypeScript + MySQL + RAG + Docker + CI

## 🧩 技能点清单

### 后端开发

| 技能 | 体现 |
|------|------|
| **Go** | Gin 框架 RESTful API，database/sql 原生 MySQL 操作 |
| **RESTful API** | 15 个接口，标准 CRUD + 搜索，统一 JSON 响应 |
| **数据库迁移** | Goose 管理 5 个迁移文件，up/down 可回滚 |
| **YAML 配置** | 分层配置（默认值 → config.yaml → config-local.yaml）|
| **中间件** | CORS、logger、recovery（Gin default） |
| **ID 生成** | 8 位随机字符串，无自增 ID 暴露 |
| **Docker** | 多阶段构建，最终镜像 ~15MB，alpine 运行 |

### AI / RAG

| 技能 | 体现 |
|------|------|
| **向量检索** | chromem-go 纯 Go 内存向量库，768 维 embedding |
| **混合搜索** | 向量相似度 + 关键词精确过滤，无匹配不返回模糊结果 |
| **LLM 接口抽象** | LLMProvider 接口，支持 OpenAI / DeepSeek / Anthropic / Ollama |
| **本地 Embedding** | Ollama + nomic-embed-text，API 不可用时自动 cosine 降级 |
| **RAG 流水线** | 检索 → 上下文构建 → LLM 生成 → 来源展示 |
| **种子知识** | 22 条专业知识，涵盖狗/猫/狐的饲养、健康、行为 |

### 前端开发

| 技能 | 体现 |
|------|------|
| **Vue 3 + TypeScript** | `<script setup>` 组合式 API，完整类型定义 |
| **组件化** | 7 个子组件，大组件从 449 行拆到 96 行 |
| **状态管理** | 响应式 ref/reactive，props + emit 数据流 |
| **API 封装** | 类型安全的 `request<T>()`，按领域拆分的 API 层 |
| **跨端** | uni-app 框架，一套代码支持 H5 + 微信小程序 |
| **弹窗系统** | 3 个统一风格的 overlay 弹出框 |
| **交互优化** | URL 编码处理、错误反馈、加载状态 |

### 工程化

| 技能 | 体现 |
|------|------|
| **CI/CD** | GitHub Actions 3 个并行 job（go vet+test, vue-tsc+vitest, docker build）|
| **测试** | 56 个后端测试 + 23 个前端测试，含性能基准 |
| **性能基准** | embedding 240ns/op，keywordMatch 0 alloc |
| **Git 规范** | conventional commit，功能粒度拆分 |
| **文档驱动** | PRD → Plan → Code 标准化流程，6 份 PRD + 6 份 Plan |
| **Docker** | 多阶段构建，从 golang:1.26-alpine → alpine:3.20 |

### AI 辅助开发（Loop Engine）

| 技能 | 体现 |
|------|------|
| **Agent 体系** | 7 个领域专用 agent，各有独立工具权限 |
| **工作流编排** | 3 个多 agent 工作流，支持并行执行 |
| **自动纠错** | 修复→验证→再修复→再验证 自动循环，最多 5 轮 |
| **知识管理** | 持久记忆体系，偏好/决策/反馈跨会话自动加载 |

## 📊 数据指标

```
后端测试: 56 个    前端测试: 23 个
Go 版本: 1.26      MySQL: 5 张表
API 接口: 15 个    知识条目: 22 条
Docker 镜像: ~15MB  CI Job: 3 个并行
组件数: 7 个子组件  迁移文件: 5 个
PRD: 6 份          Plan: 6 份
```

## 🔗 相关链接

- 项目源码: [github.com/kiripeng214/all-workspace](https://github.com/kiripeng214/all-workspace)
- CI 状态: [Actions](https://github.com/kiripeng214/all-workspace/actions)
- 演示: [pet-applet.gif](https://github.com/kiripeng214/all-workspace/raw/main/doc/pet-applet.gif)
