# 🐾 宠物助手 (Pet Applet)

[![CI](https://github.com/kiripeng214/all-workspace/actions/workflows/ci.yml/badge.svg)](https://github.com/kiripeng214/all-workspace/actions/workflows/ci.yml)

<img src="https://github.com/kiripeng214/all-workspace/raw/main/doc/pet-applet.gif" width="100%" alt="宠物助手演示">

全栈宠物管理应用 + AI 知识库。支持多宠物信息管理、喂养计划与喂养记录追踪，集成 RAG（检索增强生成）知识库，提供宠物养护智能问答。

## 功能

| 模块 | 功能 |
|------|------|
| 🐕 宠物管理 | 增删改查宠物，头像/品种/生日/体重信息 |
| 🍖 喂养计划 | 按时间/食物/分量管理每日喂养计划 |
| 📋 喂养记录 | 记录每次喂养详情，按日/按宠物查看 |
| 🧠 宠物知识库 | RAG + Ollama 本地 embedding 的智能问答，按品种自动匹配知识 |

## 项目结构

```
pet-applet/
├── backend/                # Go API 服务端
│   ├── main.go                 # 入口，Gin 路由 + CORS + 知识库初始化
│   ├── config/                 # YAML 配置 + 本地覆盖 (config-local.yaml)
│   ├── database/               # MySQL 连接 + Goose 迁移
│   ├── migrations/             # 4 个 Goose 迁移文件
│   ├── handlers/               # HTTP 处理器 (含 knowledge.go RAG 接口)
│   ├── knowledge/              # RAG 知识库 (chromem-go + LLM Provider)
│   │   ├── knowledge.go            # 向量检索 + 混合搜索 + embedding 函数
│   │   ├── llm.go                  # LLMProvider 接口 + OpenAI 实现
│   │   ├── anthropic.go            # Anthropic 实现
│   │   ├── db.go                   # MySQL 加载知识条目
│   │   └── seed.go                 # 知识条目结构定义
│   ├── models/                 # 数据模型
│   └── data/                   # 种子数据 JSON (seed_knowledge.json)
│
├── miniprogram/            # uni-app + Vue 3 + TypeScript 前端
│   ├── src/
│   │   ├── api/                  # HTTP 客户端层
│   │   │   ├── request.ts            # 通用请求封装
│   │   │   ├── knowledge.ts          # 知识库 API
│   │   │   └── pets/schedules/...
│   │   ├── pages/
│   │   │   ├── pets/                 # 宠物列表 / 详情 (含组件拆分) / 编辑
│   │   │   │   └── components/       # PetInfoCard / TodayRecords / ScheduleList / AvatarPicker
│   │   │   ├── schedules/            # 喂养计划管理 (ScheduleForm 组件)
│   │   │   │   └── components/       # ScheduleForm
│   │   │   ├── records/              # 喂养记录查看
│   │   │   └── knowledge/            # RAG 知识库页面
│   │   └── __tests__/            # 23 个测试用例
│   └── pages.json              # 路由定义
│
├── .claude/                # Claude Code Loop Engine
│   ├── commands/               # 11 个 /project:* 命令
│   ├── agents/                 # 7 个领域 agent
│   ├── skills/                 # 3 个技能
│   ├── workflows/              # 3 个工作流
│   ├── rules/                  # 后端/前端/通用编码规范
│   └── memory/                 # 持久记忆
│
├── doc/                     # 设计文档 (PRD + Plan)
│   ├── prd/                    # 产品需求文档
│   └── plan/                   # 技术方案
│
└── .github/workflows/      # CI (go vet + test, vue-tsc + vitest, docker build)
```

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.26 + [Gin](https://github.com/gin-gonic/gin) + MySQL 8.0 |
| 前端 | [uni-app 3](https://uniapp.dcloud.net.cn/) + Vue 3 + TypeScript + Vite 5 |
| 数据库迁移 | [Goose](https://github.com/pressly/goose) (4 个迁移文件) |
| 向量检索 | [chromem-go](https://github.com/philippgille/chromem-go) (纯 Go 内存向量库) |
| AI Embedding | Ollama + nomic-embed-text / 本地 cosine 降级 |
| AI 问答 | LLMProvider 接口 (OpenAI / DeepSeek / Anthropic / 降级纯检索) |
| 容器化 | Docker 多阶段构建 (最终 ~15MB) |
| CI | GitHub Actions (后端 vet+test, 前端 type-check+test, docker build) |

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0
- (可选) Ollama + nomic-embed-text — 本地 embedding

### 1. 数据库

```sql
CREATE DATABASE pet_applet CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

连接配置在 `backend/config/config.yaml`，本地覆盖写 `config-local.yaml`（不上传 Git）。

### 2. 启动后端

```bash
cd backend
go build -o pet-applet-server . && ./pet-applet-server
```

首次启动自动执行 Goose 迁移建表 + 插入种子知识数据。

### 3. 启动前端

```bash
cd miniprogram
npm install
npm run dev:h5
```

### 4. (可选) 本地 Embedding

```bash
ollama pull nomic-embed-text
```

在 `config/config-local.yaml` 配置：

```yaml
llm:
  embedding_url: "http://localhost:11434/api/embeddings"
  embedding_model: "nomic-embed-text"
```

## API 概览

| 方法 | 路径 | 说明 |
|------|------|------|
| `GET` | `/api/pets` | 宠物列表 |
| `POST` | `/api/pets` | 新增宠物 |
| `GET/PUT/DELETE` | `/api/pets/:id` | 宠物详情/更新/删除 |
| `GET/POST` | `/api/pets/schedules/:petId` | 喂养计划列表/新增 |
| `PUT/DELETE` | `/api/schedules/:id` | 喂养计划更新/删除 |
| `GET` | `/api/pets/records/:petId` | 喂养记录 |
| `GET` | `/api/pets/records/today/:petId` | 今日记录 |
| `POST` | `/api/pets/records/:petId` | 新增喂养记录 |
| `DELETE` | `/api/records/:id` | 删除记录 |
| `GET` | `/api/knowledge/search?q=&breed=` | RAG 知识库搜索 + AI 问答 |
| `GET` | `/api/meta/breeds` | 动物类型和品种列表 |

## RAG 知识库架构

```
用户输入
    │
    ▼
chromem-go 向量检索 ──→ 关键词混合过滤
    │
    ▼
LLMProvider (接口)
    ├── OpenAIProvider
    ├── AnthropicProvider
    └── 降级：纯检索结果
    │
    ▼
前端展示：AI 回答 + 参考来源
```

- **16 条种子知识**：涵盖狗/猫/通用宠物养护（疫苗接种、驱虫、饮食、品种特性等）
- **混合搜索**：向量相似度 + 关键词精确匹配，无匹配不返回模糊结果
- **自动降级**：API 不可用时自动用本地 cosine embedding，知识库不受影响

## Loop Engine

| 命令 | 说明 |
|------|------|
| `/project:loop-task` | 全流程：PRD → Plan → 实现 (顺序/并行模式) |
| `/project:component-split` | 大组件拆小组件 + 自动验证 |
| `/project:deploy-build` | 生产构建 (Docker 镜像 + 微信小程序) |
| `/project:push-code` | 检查后推送代码 |
| `/project:build-check` | 提交流前 4 步构建检查 |
| `/project:type-check` | TypeScript 类型检查 |
| `/project:loop-health` | Loop Engine 健康检查 |
| `/project:loop-engineering` | 自动审查调整 Loop Engine 体系 |

详情见 [.claude/README](.claude/)。

## 许可证

MIT
