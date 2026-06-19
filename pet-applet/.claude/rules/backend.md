---
description: Go 后端开发规范（Gin + MySQL）
---

# Go Backend 开发规范

## 项目结构

```
backend/
├── main.go         # 入口，注册路由
├── config/         # 配置加载（config.go + config.yaml）
├── database/       # 数据库连接 + 自动迁移
├── handlers/       # HTTP handler 函数
└── models/         # 数据模型结构体
```

## 编码规范

- 新增 handler 时在 `handlers/` 下创建对应文件，并在 `main.go` 中注册路由
- 表结构和模型保持同步更新（`models/models.go` + `database/database.go` 中的 `migrate()`）
- ID 统一使用 `handlers/generateID()` 生成 8 位随机小写字母+数字
- 所有 handler 返回 `gin.H` JSON 响应
- 配置优先：默认值在 `config/config.go` 的 `defaultConfig()` 中，`config.yaml` 可选覆盖
- 错误提示使用中文

## 数据库约定

- 表使用 InnoDB 引擎 + utf8mb4 字符集
- 主键使用 VARCHAR(36)，外键添加 ON DELETE CASCADE
- 时间字段使用 BIGINT 存储毫秒时间戳
- 新增字段需同步更新：DDL → 模型 struct → 所有相关 SQL 查询

## 常用命令

```bash
cd backend
go build -o pet-applet-server .
go vet ./...
```
