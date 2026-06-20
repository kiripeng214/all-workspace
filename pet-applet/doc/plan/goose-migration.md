---
status: draft
author: loop-task
created: 2026-06-20
based-on: doc/prd/goose-migration.md
---

# Plan: 替换为 Goose 迁移管理

## 1. 方案概述

- 安装 `github.com/pressly/goose/v3` 库
- 重写迁移文件为 Goose 格式
- database.go 调用 goose.UpWithContext 替换自研 migration.Run
- 删除自研 database/migration/migration.go

## 2. 涉及文件

| 文件 | 改动类型 | 改动内容 |
|------|---------|---------|
| `backend/go.mod` / `go.sum` | 修改 | 添加 goose v3 依赖 |
| `backend/migrations/` | 修改 | 3 个 .sql 文件改为 Goose 格式 |
| `backend/database/database.go` | 修改 | 替换为 goose.Up |
| `backend/database/migration/migration.go` | 删除 | 自研迁移器 |
| `doc/prd/goose-migration.md` | 新增 | PRD |
| `doc/plan/goose-migration.md` | 新增 | Plan |

## 3. 实施步骤

### 步骤 1: 安装 goose 依赖
`go get github.com/pressly/goose/v3`

### 步骤 2: 转换迁移文件
每条 SQL 加 `-- +goose Up` 头部和 `-- +goose Down` 尾部

### 步骤 3: 修改 database.go
用 `goose.UpWithContext` 替换 `migration.Run`

### 步骤 4: 删除自研迁移器
`rm -rf database/migration/`

### 步骤 5: 验证
`go vet ./...` + `go test ./...`

## 4. 风险与回退

| 风险 | 概率 | 应对 |
|------|------|------|
| goose 找不到迁移目录 | 低 | 先 rollback 到自研 runner（git checkout） |

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v0.1 | 初稿 | loop-task |
