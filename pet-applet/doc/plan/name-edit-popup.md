---
status: approved
author: Claude
created: 2026-06-20
based-on: doc/prd/name-edit-popup.md
---

# Plan: 详情页点击姓名弹出编辑框

## 1. 方案

- 仅前端改动，后端 PUT /pets/:id 已支持
- detail.vue: 姓名绑定 @click → 遮罩层弹窗 → 输入框 → 保存

## 2. 涉及文件

| 文件 | 改动 |
|------|------|
| miniprogram/src/pages/pets/detail.vue | 添加弹窗和保存逻辑 |

## 3. 技术要点

**事件**: `@click` 而非 `@tap`（H5 开发环境 @click 更可靠）

**弹窗**: 自定义遮罩层（绝对定位 view），不使用 uni-popup（减少组件依赖）

**API**: `updatePet(id, { name })` 后端支持局部更新

**刷新生效**: 保存成功后 `await loadData()` 刷新 pet 对象

## 4. 步骤

1. 姓名绑定 @click="openRename"
2. 编写自定义遮罩弹窗（圆角白色卡片+半透明蒙层）
3. 编写保存逻辑（校验 → API → 刷新）
4. 验证类型检查和运行测试

## 5. 风险

| 风险 | 应对 |
|------|------|
| 遮罩层与页面其他交互冲突 | 遮罩全屏覆盖，阻挡后续事件 |
