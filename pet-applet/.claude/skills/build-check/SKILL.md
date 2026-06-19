---
name: build-check
description: 构建检查 — 提交前验证后端编译、前端类型检查和前端编译
trigger: pre-commit
---

# Build Check

Git 提交流前自动运行，验证以下方面都能通过。

## 步骤

1. 后端编译检查：`cd backend && go build -o /dev/null ./...`
2. 前端类型检查：`cd miniprogram && npx vue-tsc --noEmit`
3. 前端构建检查：`cd miniprogram && npx vite build --mode production 2>/dev/null`（验证配置无错误）

全部通过才允许提交。任一失败则列出详细错误。

## 实施阶段的操作规范

在 `/loop-task` 的实施阶段（步骤 3）执行完修改后，必须额外执行：

- 确认 PRD 中每一条验收标准都有对应的代码实现，而非只跑编译
- 前端 UI 改动：检查 CSS 值（rpx、颜色、位置）与 PRD 一致
- 弹窗/交互类改动：检查事件绑定（@click/@tap）和显示条件（v-if）
