---
name: build-check
description: 构建检查 — 提交流前验证后端编译和前端类型检查
trigger: pre-commit
---

# Build Check

Git 提交流前自动运行，验证两方面都能通过。

## 步骤

1. 后端编译检查：`cd backend && go build -o /dev/null ./...`
2. 前端类型检查：`cd miniprogram && npx vue-tsc --noEmit`

两者全部通过才允许提交。任一失败则列出详细错误。
