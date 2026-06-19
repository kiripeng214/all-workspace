---
status: approved
author: Claude
created: 2026-06-20
based-on: doc/prd/quick-rename.md
---

# Plan: 宠物名称快速编辑

## 1. 方案概述

- 仅涉及前端，后端 `PUT /pets/:id` 已支持局部更新（只传 name 字段）
- 在详情页 `pets/detail.vue` 的宠物名字上绑定 `@tap` 事件，弹出 uni-popup 改名弹窗
- 调用 `updatePet(id, { name })` 保存，成功后 `loadData()` 刷新

## 2. 涉及文件

| 文件 | 改动类型 | 改动内容 |
|------|---------|---------|
| `miniprogram/src/pages/pets/detail.vue` | 修改 | 名字添加 @tap 事件 + 改名弹窗 + 保存逻辑 |

## 3. 实施步骤

```mermaid
graph LR
  A[名字绑 @tap] --> B[uni-popup 弹窗]
  B --> C[updatePet API]
  C --> D[loadData 刷新]
```

### 步骤 1: 名字绑定点击事件
- `<text class="name">` 添加 `@tap="openRename"`
- 名字添加虚线边框 + `:active` 态提示可点击

### 步骤 2: 改名弹窗
- 使用 `uni-popup`，`:show` 控制显隐
- 弹窗标题"修改名称"，input 预填当前名，确认/取消按钮

### 步骤 3: 保存逻辑
- 空值校验
- `updatePet(id, { name })` 保存（后端支持局部更新）
- 保存中禁用按钮防重复提交
- 成功后 `await loadData()` 刷新

## 4. 风险与回退

| 风险 | 概率 | 应对 |
|------|------|------|
| 与"编辑"按钮功能重叠 | 低 | 改名只改名字，编辑进完整表单，各司其职 |

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v1.0 | 定稿：详情页点击名字改名 | Claude |
