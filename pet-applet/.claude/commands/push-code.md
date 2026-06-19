---
name: project:push-code
description: 推送代码到远程仓库 — 先检查状态、运行检查，确认后推送
---

# /project:push-code

将本地提交推送到远程仓库。

## 执行步骤

### 1. 检查状态

```bash
cd backend && go vet ./... && go test ./...
cd miniprogram && npx vue-tsc --noEmit && npx vitest run
```

### 2. 确认推送

```bash
git push origin main
```
