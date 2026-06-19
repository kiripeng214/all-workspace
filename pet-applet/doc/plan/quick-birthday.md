---
status: approved
author: Claude
created: 2026-06-20
based-on: doc/prd/quick-birthday.md
---

# Plan: 详情页点击生日快速修改

## 1. 方案

- 仅前端细节改动
- detail.vue: 🎂标签绑@click → overlay弹窗 → picker选择 → 保存

## 2. 涉及文件

| 文件 | 改动 |
|------|------|
| miniprogram/src/pages/pets/detail.vue | 添加生日弹窗和保存逻辑 |
| miniprogram/src/api/pets.ts | 无需改动（updatePet已支持） |

## 3. 技术要点

**事件**: @click（与改名弹窗一致，已验证可用）

**弹窗**: 复用 overlay 模式（与改名弹窗结构一致）

**日期选择**: uni-app `<picker mode="date">`，与 edit.vue 用法一致

**API**: `updatePet(id, { birthday })` 后端支持局部更新

## 4. 步骤

1. 🎂标签绑 @click="openBirthdayEdit"
2. 添加 overlay 弹窗（与改名弹窗结构一致）
3. 弹窗内放 date picker + 确认/取消按钮
4. 保存逻辑（校验 → updatePet → loadData 刷新）
5. 恢复改名弹窗不影响

## 5. 风险

| 风险 | 应对 |
|------|------|
| 生日为空时显示什么 | 无生日则不显示🎂标签（v-if="pet.birthday"），不会触发点击 |
| 两个弹窗同时打开 | openBirthdayEdit 前关闭改名弹窗 |
