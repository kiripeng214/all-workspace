---
status: draft
author: loop-task
created: 2026-06-20
---

# PRD: 数据库迁移系统

## 1. 背景与动机

- 当前建表 SQL 硬编码在 Go 代码中（`database/migrate()`），修改表结构需要改代码
- 无法追踪数据库版本变更，无法回滚
- 团队协作时无法对齐数据库版本

## 2. 需求描述

### 2.1 功能列表

| # | 功能 | 优先级 | 说明 |
|---|------|--------|------|
| 1 | SQL 迁移脚本 | P0 | 每个迁移一个 `.sql` 文件，按序号执行 |
| 2 | 自动追踪 | P0 | `schema_migrations` 表记录已执行的迁移 |
| 3 | 启动时自动执行 | P0 | 未执行的迁移按顺序自动执行 |
| 4 | 保留 GO 内联回退 | P1 | 迁移失败时保留原有 `migrate()` 作为降级 |

### 2.2 目录结构

```
backend/
├── migrations/
│   ├── 001_create_pets.sql
│   ├── 002_create_feeding_schedules.sql
│   └── 003_create_feeding_records.sql
└── database/
    └── migration/
        └── migration.go    # 迁移执行器
```

### 2.3 用户流程

```
启动后端
  → 连接数据库
  → 检查 schema_migrations 表（不存在则创建）
  → 读取 migrations/ 目录下所有 .sql 文件
  → 按文件名排序，逐个执行未运行过的迁移
  → 执行成功后在 schema_migrations 记录文件名
  → 启动 HTTP 服务
```

## 3. 验收标准

- [ ] `migrations/` 目录下 3 个 sql 文件，分别对应 pets, feeding_schedules, feeding_records
- [ ] 启动后自动创建 `schema_migrations` 表并执行全部迁移
- [ ] 再次启动不重复执行已完成的迁移
- [ ] 迁移失败时打印错误并退出
- [ ] `go vet ./...` 通过
- [ ] `go test ./...` 通过

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v0.1 | 初稿 | loop-task |
