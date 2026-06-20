---
status: draft
author: loop-task
created: 2026-06-20
---

# Plan: 知识库存入 MySQL

## 1. 方案概述

- 新增 `knowledge_entries` 表 + 迁移文件
- seed 数据改为 JSON 文件 + 迁移 INSERT
- 启动时 knowledge.Init 从 MySQL 加载 → 构建 chromem-go 索引
- 移除 seed.go 硬编码数据（保留 struct 定义）

## 2. 涉及文件

| 文件 | 改动类型 | 改动内容 |
|------|---------|---------|
| `backend/migrations/004_create_knowledge_entries.sql` | 新增 | 建表 + 种子数据 INSERT |
| `backend/data/seed_knowledge.json` | 新增 | 15 条种子数据 JSON |
| `backend/knowledge/db.go` | 新增 | 从 MySQL 加载知识条目 |
| `backend/knowledge/knowledge.go` | 修改 | Init 接受 *sql.DB，从 DB 加载 |
| `backend/knowledge/seed.go` | 修改 | 只保留 KnowledgeEntry 定义 |
| `backend/main.go` | 修改 | 传递 database.DB 给 knowledge.Init |

## 3. 实施步骤

### 步骤 1: 迁移文件 + 种子 JSON
- 004 migration: CREATE TABLE + INSERT 初识数据
- JSON 文件: 15 条知识条目

### 步骤 2: db.go — 从 MySQL 加载
- Query all rows → 填充 chromem-go collection

### 步骤 3: 更新调用链
- knowledge.Init 接收 *sql.DB → 从 DB 加载 → 构建向量索引

## 4. 风险与回退

| 风险 | 概率 | 应对 |
|------|------|------|
| 表不存在 | 低 | 迁移先执行，后初始化知识库 |

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v0.1 | 初稿 | loop-task |
