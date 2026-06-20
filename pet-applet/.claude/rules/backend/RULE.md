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
- 表结构和模型保持同步更新（`models/models.go` + Goose 迁移文件）
- ID 统一使用 `handlers/generateID()` 生成 8 位随机小写字母+数字
- 所有 handler 返回 `gin.H` JSON 响应
- 配置优先：默认值在 `config/config.go` 的 `defaultConfig()` 中，`config.yaml` 可选覆盖
- 错误提示使用中文

## 数据库约定

- 表使用 InnoDB 引擎 + utf8mb4 字符集
- 主键使用 VARCHAR(36)，外键添加 ON DELETE CASCADE
- 时间字段使用 BIGINT 存储毫秒时间戳
- 新增字段需同步更新：DDL → 模型 struct → 所有相关 SQL 查询

## 数据库迁移（Goose）

所有表结构变更必须通过 Goose 迁移，禁止直接改表或手动执行 SQL。

### 迁移文件规范

```
backend/migrations/
├── YYYYMMDDHHMMSS_description.sql    ← 新建用时间戳
```

格式要求：
- 文件名：`YYYYMMDDHHMMSS_description.sql`（Goose 标准时间戳）
- 文件头尾必须有 Goose 注释标记：

```sql
-- +goose Up
CREATE TABLE ...;

-- +goose Down
DROP TABLE IF EXISTS ...;
```

### 编写规范

- `-- +goose Up` 和 `-- +goose Down` 之间空一行
- 每条迁移只做一件事（一张表或一个字段变更）
- Down 迁移必须能回滚 Up 的所有操作
- 有外键依赖时，Down 按依赖逆序 DROP
- 迁移文件写好先本地执行 `goose up` 验证
- 新表需同步创建模型 struct（models/）和测试

## 常用命令

```bash
cd backend
go build -o pet-applet-server .
go vet ./...
```
