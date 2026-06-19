---
status: approved
author: Claude
created: 2026-06-20
based-on: doc/prd/quick-rename.md
---

# Plan: 宠物名称快速编辑

## 1. 方案概述

- 仅涉及前端，后端 `PUT /pets/:id` 已支持局部更新（只传 name 即可）
- 在列表页 `pets/index.vue` 增加长按事件，弹出 uni-popup 改名弹窗
- 调用已有 `updatePet(id, { name })` API

## 2. 涉及文件

| 文件 | 改动类型 | 改动内容 |
|------|---------|---------|
| `miniprogram/src/pages/pets/index.vue` | 修改 | 添加长按事件、改名弹窗、保存逻辑 |

## 3. 实施步骤

```mermaid
graph LR
  A[添加长按事件] --> B[弹窗 UI]
  B --> C[保存逻辑]
  C --> D[刷新列表]
```

### 步骤 1: 添加长按事件
- 在卡片上绑定 `@longpress` 事件，传入当前 pet
- 事件触发时记录当前编辑的 pet，显示弹窗

### 步骤 2: 弹窗 UI
- 使用 uni-app 的 `uni-popup` 组件（项目已引入，参照 detail.vue 用法）
- 弹窗内容：标题「修改名称」+ input 输入框（v-model 绑定，当前名称预填）+ 确认/取消按钮
- 与现有 detail.vue 的记录喂养弹窗风格一致

### 步骤 3: 保存逻辑
- 点击确认后校验名称不为空
- 调用 `updatePet(petId, { name: newName })` 保存
- 保存期间显示 loading，按钮禁用
- 成功后关闭弹窗，调用 `loadPets()` 刷新列表

## 4. 风险与回退

| 风险 | 概率 | 应对 |
|------|------|------|
| uni-popup 组件需要导入 | 低 | 查看项目中是否已全局注册，否的话需注册 |
| 长按事件在微信小程序兼容性 | 低 | uni-app 封装了 `@longpress`，跨端兼容 |

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v0.1 | 初稿 | Claude |
