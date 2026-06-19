---
status: approved
author: Claude
created: 2026-06-20
based-on: doc/prd/detail-rename.md
---

# Plan: 详情页点击名字快速改姓名

## 1. 方案概述

- 仅涉及前端，后端 `PUT /pets/:id` 已支持局部更新
- 在详情页 `pets/detail.vue` 的宠物姓名上绑定 `@tap` → 弹 uni-popup → 修改 → 保存

## 2. 涉及文件

| 文件 | 改动类型 | 改动内容 |
|------|---------|---------|
| `miniprogram/src/pages/pets/detail.vue` | 修改 | 姓名添加 @tap 事件 + 改名弹窗 + 保存逻辑 |

## 3. 实施步骤

### 步骤 1: 姓名绑定点击事件
- `<text class="name">` 添加 `@tap="openRename"`
- 添加虚线边框 + `:active` 样式提示可点击

### 步骤 2: 改名弹窗
- 使用 `uni-popup`，`:show` 控制显隐
- 标题"修改姓名"，input 预填当前姓名，确认/取消按钮

### 步骤 3: 保存逻辑
- `updatePet(id, { name })` 保存
- 保存中禁用按钮，成功后 `loadData()` 刷新

## 4. 风险与回退

| 风险 | 概率 | 应对 |
|------|------|------|
| 点击姓名与编辑按钮重复 | 低 | 改名只改姓名，编辑进完整表单，各司其职 |

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v1.0 | 初稿 | Claude |
