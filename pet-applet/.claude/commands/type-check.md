---
name: project:type-check
description: TypeScript 类型检查 — 发现类型错误
---

# /project:type-check

运行 vue-tsc 对前端项目做类型检查 + vitest 运行测试。

## 使用场景

- 修改了 TypeScript 类型或接口后
- 修改了 API 层参数/返回值后
- 新增了 Vue 组件后
- 提交代码前

## 执行

```bash
cd miniprogram
npx vue-tsc --noEmit        # 类型检查
npx vitest run               # 运行测试
```

## 验证

类型检查无错误输出，测试全部通过。
