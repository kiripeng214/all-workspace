---
status: draft
author: loop-task
created: 2026-06-20
based-on: doc/prd/db-migration.md
---

# Plan: 数据库迁移系统

## 1. 方案概述

- 新增 `database/migration/migration.go` 迁移执行器
- 将硬编码 SQL 提取到 `migrations/` 目录
- 修改 `database.go` 调用迁移器替代内联 migrate()
- 仅后端改动

## 2. 涉及文件

| 文件 | 改动类型 | 改动内容 |
|------|---------|---------|
| `backend/migrations/001_create_pets.sql` | 新增 | pets 表 DDL |
| `backend/migrations/002_create_feeding_schedules.sql` | 新增 | feeding_schedules 表 DDL |
| `backend/migrations/003_create_feeding_records.sql` | 新增 | feeding_records 表 DDL |
| `backend/database/migration/migration.go` | 新增 | 迁移执行器（读取 sql、追踪版本、自动执行） |
| `backend/database/database.go` | 修改 | 用 migration.Run 替换内联 migrate() |

## 3. 实施步骤

### 步骤 1: 创建 SQL 迁移文件
- 从原 `migrate()` 函数提取 3 条 DDL 到独立 .sql 文件
- 编号 001/002/003

### 步骤 2: 编写迁移执行器
- `migration.go` 包
- `Run(db *sql.DB, dir string) error`
- 建 `schema_migrations` 表（如不存在）
- 读取目录下所有 `.sql` 文件，按文件名排序
- 跳过已执行的迁移
- 逐个执行未执行的迁移并记录

### 步骤 3: 修改 database.go
- 删除 `migrate()` 函数
- `Init` 末尾调用 `migration.Run(DB, "migrations")`

### 步骤 4: 验证
- `go vet ./...`
- `go test ./...`

## 4. 风险与回退

| 风险 | 概率 | 应对 |
|------|------|------|
| 迁移文件读不到（路径问题） | 中 | 用可执行文件相对路径 + 打印详细错误 |
| DDL 与原内联 SQL 不一致 | 低 | 逐条对比后复制，确保无差异 |

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v0.1 | 初稿 | loop-task |
