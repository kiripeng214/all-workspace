---
name: project:build-check
description: 构建检查 — 提交前验证后端编译、前端类型检查和前端编译
---

# /project:build-check

提交前构建检查，确保代码无编译和类型错误。

## 执行步骤

### 1. 后端编译检查

```bash
cd backend && go build -o /dev/null ./...
```

### 2. 前端类型检查

```bash
cd miniprogram && npx vue-tsc --noEmit
```

### 3. 前端测试

```bash
cd miniprogram && npx vitest run
```

### 4. 前端构建检查

```bash
cd miniprogram && npx vite build --mode production 2>/dev/null
```

## 验证

全部通过即为绿色输出。任一失败则列出详细错误信息。
