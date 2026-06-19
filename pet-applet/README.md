# 🐾 宠物助手 (Pet Applet)

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
├── .claude/              # Claude Code 配置（规则 / 技能 / 代理）
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

## 开发规范

参见 `.claude/` 目录下的规则文件：

- [通用规范](.claude/rules/general/RULE.md)
- [后端规范](.claude/rules/backend/RULE.md)
- [前端规范](.claude/rules/frontend/RULE.md)

## 许可证

MIT
