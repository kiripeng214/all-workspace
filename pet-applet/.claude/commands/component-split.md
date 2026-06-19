---
name: project:component-split
description: 前端大组件拆小组件 — 识别大组件并拆分为可复用子组件
---

# /project:component-split

将过大的 Vue 单文件组件拆分为职责单一的小组件。

## 检测标准

| 指标 | 阈值 |
|------|------|
| 文件行数 | > 250 行 |
| Template 节点 | > 50 个 |
| `<script>` 块 | > 150 行 |
| `v-if` 分支 | > 5 组 |

## 当前项目状况

```
pets/detail.vue         449 行 ← 优先拆分
schedules/index.vue     203 行
pets/edit.vue           201 行
pets/index.vue          135 行
records/index.vue        94 行
```

## 拆分步骤

1. **分析** — 识别 template 中可以独立成块的区域
2. **抽离** — 创建 `components/` 子目录，提取子组件（props + emit）
3. **替换** — 父组件中引入子组件，替换原内联代码
4. **验证** — `npx vue-tsc --noEmit && npx vitest run`

> 详情见 `.claude/skills/component-split/SKILL.md`
