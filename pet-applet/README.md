# 🐾 宠物助手 (Pet Applet)

[![CI](https://github.com/kiripeng214/all-workspace/actions/workflows/ci.yml/badge.svg)](https://github.com/kiripeng214/all-workspace/actions/workflows/ci.yml)

<video src="https://github.com/kiripeng214/all-workspace/raw/main/doc/pet-applet.mov" width="100%" controls></video>

一个宠物管理小程序，支持记录宠物信息、管理喂养计划和喂养记录。

## 项目结构

```
pet-applet/
├── backend/              # Go API 服务端
│   ├── main.go               # 入口，Gin 路由注册
│   ├── config/               # YAML 配置加载 (config.go + config.yaml)
│   ├── database/             # MySQL 连接 + 自动迁移
│   ├── handlers/             # HTTP 处理器 (pets, schedules, records, meta)
│   └── models/               # 数据模型 (Pet, FeedingSchedule, FeedingRecord)
│
├── miniprogram/          # uni-app + Vue 3 + TypeScript 前端
│   ├── src/
│   │   ├── api/              # 按领域拆分的 HTTP 客户端
│   │   ├── config/           # API_BASE_URL 配置
│   │   ├── pages/            # 页面组件
│   │   │   ├── pets/         # 列表 / 详情 / 编辑
│   │   │   ├── schedules/    # 喂养计划管理
│   │   │   ├── records/      # 喂养记录查看
│   │   │   └── index/        # 重定向至 pets/index
│   │   ├── App.vue
│   │   └── main.ts
│   ├── pages.json            # 路由定义
│   └── vite.config.ts        # Vite + uni-app 构建配置
│
├── .claude/              # Claude Code Loop Engine 体系
│   ├── settings.json          # 总控配置（目录映射 / hooks / 偏好）
│   ├── agents/               # 7 个专用子 agent
│   ├── commands/             # 11 个 /project:* 快捷命令
│   ├── skills/               # 3 个自动化技能脚本
│   ├── workflows/            # 3 个多 agent 工作流
│   ├── rules/                # 3 个领域编码规范
│   └── memory/               # 持久记忆（偏好 / 决策 / 反馈）
└── README.md
```

## 技术栈

| 层级   | 技术                                                              |
| ------ | ----------------------------------------------------------------- |
| 后端   | Go + [Gin](https://github.com/gin-gonic/gin) + MySQL              |
| 前端   | [uni-app](https://uniapp.dcloud.net.cn/) + Vue 3 + TypeScript     |
| 数据库 | MySQL 8.0 (InnoDB + utf8mb4)                                      |

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- MySQL 8.0

### 1. 数据库

创建数据库：

```sql
CREATE DATABASE pet_applet CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

默认连接配置在 `backend/config/config.yaml`，可自行修改：

```yaml
db:
  host: localhost
  port: "3306"
  user: root
  password: "123456"
  name: pet_applet
```

> 注：表结构会在后端启动时自动创建（auto-migration）。

### 2. 启动后端

```bash
cd backend
go build -o pet-applet-server .
./pet-applet-server
```

API 服务默认监听 `http://localhost:3000`。

### 3. 启动前端

```bash
cd miniprogram
npm install
npm run dev:h5        # H5 浏览器开发
# 或
npm run dev:mp-weixin  # 微信小程序开发
```

### 完整命令速查

```bash
# 后端
cd backend && go build -o pet-applet-server . && ./pet-applet-server

# 前端 H5
cd miniprogram && npm run dev:h5

# 前端 类型检查
npm run type-check

# 前端 生产构建
npm run build:h5
```

## API 概览

所有接口以 `/api` 为前缀。

| 方法     | 路径                          | 说明                   |
| -------- | ----------------------------- | ---------------------- |
| `GET`    | `/api/pets`                   | 获取宠物列表            |
| `GET`    | `/api/pets/:id`               | 获取宠物详情            |
| `POST`   | `/api/pets`                   | 新增宠物                |
| `PUT`    | `/api/pets/:id`               | 更新宠物信息            |
| `DELETE` | `/api/pets/:id`               | 删除宠物                |
| `GET`    | `/api/pets/schedules/:petId`  | 获取某宠物的喂养计划     |
| `POST`   | `/api/pets/schedules/:petId`  | 新增喂养计划            |
| `PUT`    | `/api/schedules/:id`          | 更新喂养计划            |
| `DELETE` | `/api/schedules/:id`          | 删除喂养计划            |
| `GET`    | `/api/pets/records/:petId`    | 获取某宠物的喂养记录     |
| `GET`    | `/api/pets/records/today/:petId` | 获取今日喂养记录      |
| `POST`   | `/api/pets/records/:petId`    | 新增喂养记录            |
| `DELETE` | `/api/records/:id`            | 删除喂养记录            |
| `GET`    | `/api/meta/breeds`            | 获取动物类型和品种列表   |

## Loop Engine（`.claude/` 体系）

Loop Engine 是项目的智能协作体系，让 Claude 理解项目结构、编码规范和历史决策，减少重复沟通。

### 组成

| 目录 | 数量 | 作用 |
|------|------|------|
| `agents/` | 7 | 领域专用子 agent（Go 后端、uni-app 前端、数据库、审查等），各有独立工具权限 |
| `commands/` | 11 | `/:project:*` 快捷命令，覆盖开发/测试/构建/部署全流程 |
| `skills/` | 3 | 自动化脚本，`build-check` 绑定 pre-commit 自动校验 |
| `workflows/` | 3 | 多 agent 编排工作流，支持并行任务 + 自动合并 |
| `rules/` | 3 | 编码规范（通用/后端/前端），每次对话自动加载 |
| `memory/` | — | 持久记忆，记录用户偏好、历史决策、反馈，跨会话自动加载 |

### 核心工作流

```
你输入需求
    │
    ▼
/project:loop-task   →   PRD → Plan → 实现
    │                     ├── 顺序模式：单一领域
    │                     └── 并行模式：跨领域（数据库+后端+前端并行）
    │
    ├── /component-split    →  分析大组件 → 拆分子组件 → 自动验证
    ├── /code-review        →  多维度代码审查
    ├── /deploy-build       →  后端 Docker 镜像 + 前端小程序打包
    └── /push-code          →  检查→推送
```

### 自动纠错机制

测试失败或审查发现问题时，自动进入纠错循环：

```
修复 → 重跑验证 → 还失败 → 再修复 → 再验证 → ... → 最多 5 轮
```

不中断、不等待用户确认，通过后自动继续流程。

### 快捷命令一览

| 命令 | 说明 |
|------|------|
| `/project:backend-start` | 构建并启动 Go 后端 |
| `/project:frontend-dev` | 启动前端 H5 开发服务器 |
| `/project:db-init` | 初始化 MySQL 数据库 |
| `/project:build-check` | 提交前构建检查（4 步） |
| `/project:type-check` | TypeScript 类型检查 |
| `/project:component-split` | 大组件拆小组件 |
| `/project:deploy-build` | 生产构建（后端镜像 + 前端小程序） |
| `/project:push-code` | 检查后推送代码 |
| `/project:loop-task` | 全流程：PRD→Plan→实现 |
| `/project:loop-health` | Loop Engine 健康检查 |
| `/project:loop-engineering` | 自动审查并调整 Loop Engine 体系 |

## 许可证

MIT
