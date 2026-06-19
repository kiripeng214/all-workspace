---
status: approved
author: Claude
created: 2026-06-20
based-on: doc/prd/quick-rename.md
---

# Plan: 宠物名称快速编辑

## 1. 方案概述

- 仅涉及前端，后端 `PUT /pets/:id` 已支持局部更新
- 在详情页 `pets/detail.vue` 的宠物名字上绑定 `@tap` 事件，弹出改名弹窗
- 调用已有 `updatePet(id, { name })` API，保存后刷新页面

## 2. 涉及文件

| 文件 | 改动类型 | 改动内容 |
|------|---------|---------|
| `miniprogram/src/pages/pets/detail.vue` | 修改 | 名字添加 @tap 事件 + 改名弹窗 + 保存逻辑 |

## 3. 实施步骤

```mermaid
graph LR
  A[名字绑 tap 事件] --> B[改名弹窗]
  B --> C[保存逻辑 + 刷新]
```

### 步骤 1: 名字绑定点击事件
- `<text class="name">` 添加 `@tap="openRename"`
- 与现有"编辑"按钮并存：改名是轻量操作，编辑按钮进完整表单

### 步骤 2: 改名弹窗
- 使用 `uni-popup`（项目已用，与"记录喂养"弹窗一致）
- 弹窗标题"修改名称"，input 预填当前名，确认/取消按钮
- `:show` 控制显隐

### 步骤 3: 保存逻辑
- 空值校验，`updatePet(id, { name })` 保存
- 保存中禁用按钮，成功后 `loadData()` 刷新
- 与前端已有 `updatePet` 类型对齐：`UpdatePetParams` 是 `Partial<CreatePetParams>`，只传 `name` 即可

## 4. 风险与回退

| 风险 | 概率 | 应对 |
|------|------|------|
| `uni-popup` 已使用 | 无风险 | 沿用现有弹窗模式 |
| 名字 tap 与编辑按钮重复 | 低 | 改名只改名字，编辑进完整表单，功能不同 |

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v0.1 | 初稿（列表页长按方案） | Claude |
| 2026-06-20 | v0.2 | 修正为详情页点击名字 | Claude |
