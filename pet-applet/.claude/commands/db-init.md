---
name: project:db-init
description: 初始化 MySQL 数据库和表
---

# /project:db-init

创建宠物助手所需的 MySQL 数据库及表结构。

## 执行步骤

1. 登录 MySQL：
   ```bash
   mysql -u root -p
   ```

2. 创建数据库（如果不存在）：
   ```sql
   CREATE DATABASE IF NOT EXISTS pet_applet CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

3. 项目启动时会自动建表（`database/migrate()`），无需手动执行 DDL。

## 当前表结构

- **pets** — 宠物信息（id, avatar, name, breed, birthday, weight, notes, created_at）
- **feeding_schedules** — 喂养计划（id, pet_id, time, food_type, amount）
- **feeding_records** — 喂养记录（id, pet_id, schedule_id, time, food_type, amount, notes, created_at）

所有表使用 InnoDB 引擎 + utf8mb4 字符集。
