---
name: type-check
description: TypeScript 类型检查 — 发现类型错误
---

# TypeScript Type Check

运行 vue-tsc 对前端项目做类型检查（不生成输出文件）。

## 使用场景

- 修改了 TypeScript 类型或接口后
- 修改了 API 层参数/返回值后
- 新增了 Vue 组件后

## 执行

```bash
cd miniprogram && npx vue-tsc --noEmit
```
