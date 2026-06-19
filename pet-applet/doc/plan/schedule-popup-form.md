---
status: draft
author: loop-task
created: 2026-06-20
based-on: doc/prd/schedule-popup-form.md
---

# Plan: 喂养计划表单改为弹出框

## 1. 方案概述

- 将 `ScheduleForm.vue` 的内联渲染改为弹出框模式
- 在前端领域内完成，不涉及后端和数据库
- `schedules/index.vue` 新增 `showForm` 控制弹出框显示
- `ScheduleForm.vue` 内容包裹在遮罩层中，实现弹出框样式

## 2. 涉及文件

| 文件 | 改动类型 | 改动内容 |
|------|---------|---------|
| `miniprogram/src/pages/schedules/index.vue` | 修改 | 新增 `showForm` 控制状态、添加 FAB 按钮、传递 `showForm` prop |
| `miniprogram/src/pages/schedules/components/ScheduleForm.vue` | 修改 | 添加遮罩层包裹、接收 `show` prop 控制显隐、调整样式为弹出框 |
| `miniprogram/src/__tests__/components/schedules.test.ts` | 新增 | 添加弹出框交互测试 |

## 3. 实施步骤

### 步骤 1: 修改 ScheduleForm.vue
- template 用 `<view v-if="show" class="overlay">` 包裹现有内容
- 新增 `show` prop
- 新增 `close` emit（点击遮罩层时触发）
- 点击取消按钮 emit `cancel`
- 样式：遮罩层 fixed 全屏半透明黑色、表单卡片白色居中圆角

### 步骤 2: 修改 schedules/index.vue
- 新增 `showForm` ref（默认 false）
- 新增 "添加计划" 按钮代替原有的内联渲染
- 编辑计划时：设置 `editingSchedule` + `showForm = true`
- 添加计划时：清空 `editingSchedule` + `showForm = true`
- 表单取消/关闭时：`showForm = false`
- 移除原有的 `<ScheduleForm>` 内联渲染，改为条件渲染

### 步骤 3: 写测试
- 测试弹出框显隐
- 测试添加弹出表单空白
- 测试编辑弹出表单预填

### 步骤 4: 验证
- `vue-tsc --noEmit`
- `vitest run`

## 4. 风险与回退

| 风险 | 概率 | 应对 |
|------|------|------|
| 遮罩层在微信小程序中滚动穿透 | 低 | 遮罩层加 `@touchmove.prevent` |
| picker 组件在弹出框中行为异常 | 低 | 保留现有 picker 用法，H5 和微信均原生支持 |

---

## 变更记录

| 日期 | 版本 | 变更内容 | 作者 |
|------|------|---------|------|
| 2026-06-20 | v0.1 | 初稿 | loop-task |
