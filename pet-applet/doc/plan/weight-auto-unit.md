---
status: draft
author: loop-task
created: 2026-06-20
based-on: doc/prd/weight-auto-unit.md
---

# Plan: 体重字段自动拼接单位

## 1. 方案概述

- 前端 `edit.vue` 修改：替换体重文本输入框为数字输入 + 单位选择器
- 纯前端改动，不涉及后端和数据库

## 2. 涉及文件

| 文件 | 改动类型 | 改动内容 |
|------|---------|---------|
| `miniprogram/src/pages/pets/edit.vue` | 修改 | template 拆为数字+单位；script 增加解析/拼接逻辑；style 增加布局样式 |

## 3. 实施步骤

### 步骤 1: 修改 template
- 替换 `.input` 为 `.weight-row` 容器
- 左：数字输入框 `<input type="digit" v-model="weightValue">`
- 右：单位 picker `<picker :range="['kg', 'g']">`

### 步骤 2: 修改 script
- 新增 `weightValue` / `weightUnit` / `weightUnits` 响应式数据
- `onWeightUnitChange` — picker 回调更新单位
- `onLoad` 编辑模式 — 正则解析 `"15kg"` 为数字+单位
- `onSubmit` — 拼接 `weightValue + weightUnit` → `form.weight`

### 步骤 3: 验证
- `vue-tsc --noEmit`
- `vitest run`

## 4. 风险与回退

| 风险 | 概率 | 应对 |
|------|------|------|
| 后端有存量数据格式不统一 | 低 | 正则 `/^([\d.]+)(kg\|g)$/` 匹配不到时保持原值 |

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v0.1 | 初稿 | loop-task |
