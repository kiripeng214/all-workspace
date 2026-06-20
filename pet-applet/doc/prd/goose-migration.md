---
status: draft
author: loop-task
created: 2026-06-20
---

# PRD: 替换为 Goose 迁移管理

## 1. 背景与动机

- 之前自研的迁移执行器功能简陋，不支持 down 回滚、不支持 CLI 操作
- Goose 是 Go 生态最成熟的迁移工具，支持 up/down/status CLI 命令 + 库调用
- 社区标准，团队协作零学习成本

## 2. 需求描述

### 2.1 功能列表

| # | 功能 | 优先级 | 说明 |
|---|------|--------|------|
| 1 | 迁移文件改为 Goose 格式 | P0 | `-- +goose Up` / `-- +goose Down` 注释标记 |
| 2 | 启动时自动执行迁移 | P0 | database.Init 调用 goose.Up 代替自研 runner |
| 3 | CLI 管理 | P0 | 可用 `goose up` / `goose down` / `goose status` |
| 4 | 移除自研迁移器 | P1 | 删除 database/migration/migration.go |

### 2.2 迁移文件格式

```sql
-- +goose Up
CREATE TABLE ...;

-- +goose Down
DROP TABLE IF EXISTS ...;
```

## 3. 验收标准

- [ ] 迁移文件带 `-- +goose Up/Down` 标记
- [ ] 启动时自动执行全部未跑迁移
- [ ] 再次启动不重复执行
- [ ] `goose status` 可查看迁移状态
- [ ] `go vet ./...` 通过
- [ ] `go test ./...` 通过

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v0.1 | 初稿 | loop-task |
