---
name: component-split
description: 前端大组件拆小组件 — 识别大组件并拆分为可复用子组件
trigger: manual
---

# Component Split

将过大的 Vue / uni-app 单文件组件拆分为职责单一的小组件，提升可维护性和复用性。

## 检测标准

遇到以下情况应考虑拆分：

| 指标 | 阈值 | 说明 |
|------|------|------|
| 行数 | > 250 行 | 超出单屏幕可浏览范围 |
| template 内节点 | > 50 个 DOM 节点 | 逻辑块过多，难以追踪渲染路径 |
| `<script>` 块 | > 150 行 | 承载了过多的响应式数据和逻辑 |
| `v-if/v-else` 分支 | > 5 组 | 组件在不同状态下渲染完全不同的内容 |
| `computed` + `watch` | > 8 个 | 派生逻辑过于集中 |
| template 注释分组 | 存在 "== 区块 == " 注释 | 作者自己都在用注释做视觉分割 |

## 拆分策略

### 1. 识别可拆分区块

在大组件中找到以下典型模式：

```
<!-- 宠物信息卡片 -->       → 抽为 PetInfoCard.vue
<!-- 今日喂养记录 -->         → 抽为 TodayRecords.vue
<!-- 喂养计划列表 -->          → 抽为 ScheduleList.vue
<!-- 快捷记录表单 -->          → 抽为 QuickRecordForm.vue
<!-- 弹窗/模态框内容 -->       → 抽为 xxxModal.vue / xxxPopup.vue
```

### 2. 提取子组件步骤

**步骤 A：建立组件目录**

```
src/pages/pets/
├── detail.vue               # 父组件，只保留布局和编排逻辑
├── components/              # 新建目录
│   ├── PetInfoCard.vue
│   ├── TodayRecords.vue
│   ├── ScheduleList.vue
│   └── QuickRecordForm.vue
```

**步骤 B：确定组件接口**

| 传参方式 | 用途 |
|----------|------|
| `props` | 父→子 数据传递 |
| `emit` | 子→父 事件通知 |
| `v-model` | 双向绑定（表单类组件） |

提取原则：
- **props 扁平化** — 避免传递整个对象，只传子组件真正使用的字段
- **事件名用 dash-case** — 如 `@record-added="onRecordAdded"`
- **不跨两级传递** — 子组件不应透传 props 给孙组件（用 slot 或重构）

**步骤 C：生成子组件骨架**

```vue
<script setup lang="ts">
// 1. 从父组件复制相关 props
// 2. 复制相关数据和计算属性
// 3. 复制相关方法（emit 替换直接调用）
// 4. 添加 defineEmits 声明
</script>

<template>
  <view class="component-name">
    <!-- 从父组件复制相关 template -->
  </view>
</template>

<style scoped>
/* 复制相关样式 */
</style>
```

**步骤 D：父组件替换**

```diff
- <!-- 宠物信息卡片 -->
- <view class="info-card">
-   ...
- </view>
+ <PetInfoCard :pet="pet" @edit="goEdit" />
```

## 拆分验证（自动执行）

拆分完成后**必须自动执行以下验证**，不通过则进入纠错循环（同 loop-task 机制）：

### 验证步骤

```bash
cd miniprogram
npx vue-tsc --noEmit      # 类型检查
npx vitest run              # 测试通过
```

### 自动纠错循环

```
❌ vue-tsc 失败 → 分析报错 → 修复类型错误 → 重跑
❌ vitest 失败 → 分析报错 → 修复逻辑错误 → 重跑
✅ 全部通过 → 继续下一步
```

最大 3 轮，同一错误两次则终止。

### 检查清单

拆分前确认：

- [ ] 子组件 `:props` 类型已定义（TypeScript interface）
- [ ] 子组件 `defineEmits` 已声明
- [ ] 子组件样式用 `scoped` 隔离（`<style scoped>`）
- [ ] 父组件 import 路径正确
- [ ] 子组件不依赖父组件的 `provide/inject`（如无必要）
- [ ] 子组件不直接调用父组件函数（用 emit）

## 项目当前状况

```
pets/detail.vue          96 行 ← 已拆分（原 449 行）
  ├── PetInfoCard.vue    275 行 ← 宠物信息卡片 + 重命名/生日弹窗
  ├── TodayRecords.vue   148 行 ← 今日喂养记录 + 快捷记录弹窗
  └── ScheduleList.vue    71 行 ← 喂养计划列表
pets/edit.vue           167 行 ← 已拆分（原 201 行）
  └── AvatarPicker.vue    46 行 ← 表情网格选择器
schedules/index.vue     124 行 ← 已拆分（原 203 行）
  └── ScheduleForm.vue   116 行 ← 内联添加/编辑表单
records/index.vue        94 行
pets/index.vue          135 行
```

**已完成 3 个组件拆分：**

| 父组件 | 原始 | 精简后 | 子组件 |
|--------|------|--------|--------|
| `detail.vue` | 449 | **96** | `PetInfoCard`, `TodayRecords`, `ScheduleList` |
| `edit.vue` | 201 | **167** | `AvatarPicker` |
| `schedules/index.vue` | 203 | **124** | `ScheduleForm` |

剩余候选：`records/index.vue`（94 行）和 `pets/index.vue`（135 行）行数尚可，暂无需拆分。

## 使用方式

```bash
# 审查 detail.vue 的拆分方案
请分析 src/pages/pets/detail.vue 适合拆分成哪些子组件

# 执行拆分
请将 detail.vue 的宠物信息卡片部分提取为 PetInfoCard.vue

# 验证
npx vue-tsc --noEmit && npx vitest run
```
