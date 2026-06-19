---
description: 全局行为规范
---

# 全局行为规范

## 沟通语言

- 与用户交流、错误提示、注释使用中文
- 变量名、函数名、类型名使用英文

## 变更规范

- 修改后端时同步检查前端 API 调用是否匹配
- 修改表结构时同步更新 DDL、模型、所有 SQL 查询
- 新增页面时同步更新 `pages.json` 路由注册

## 代码风格

- TypeScript 接口名称使用 PascalCase（如 `FeedingRecord`）
- Go 函数名使用 PascalCase，变量使用 camelCase
- Vue 组件使用 `<script setup lang="ts">` 组合式 API
- 配置优先原则：提供合理默认值，外部配置可选覆盖
