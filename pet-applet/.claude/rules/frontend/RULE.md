---
description: uni-app 前端开发规范（Vue 3 + TypeScript）
---

# Frontend 开发规范

## 项目结构

```
miniprogram/
├── src/
│   ├── api/          # 按领域拆分的 API 客户端
│   │   ├── request.ts   # 通用 HTTP 请求封装
│   │   ├── pets.ts      # 宠物相关接口
│   │   ├── schedules.ts # 喂养计划接口
│   │   ├── records.ts   # 喂养记录接口
│   │   └── meta.ts      # 元数据接口
│   ├── config/       # API_BASE_URL 等配置
│   └── pages/        # 页面组件
│       ├── pets/         # 列表/详情/编辑
│       ├── schedules/    # 喂养计划管理
│       └── records/      # 喂养记录查看
├── pages.json        # 路由定义
└── vite.config.ts    # Vite + uni-app 构建配置
```

## 编码规范

- 页面路径在 `pages.json` 中定义，源码按 `src/pages/xxx/` 组织
- API 调用封装在 `src/api/` 下，按领域拆分为独立文件
- 统一使用 `src/api/request.ts` 的 `request<T>()` 进行 HTTP 请求
- `API_BASE_URL` 在 `src/config/index.ts` 中配置，默认 `http://localhost:3000/api`
- 使用 uni-app 内置 API（`uni.navigateTo`, `uni.showToast`, `uni.showModal` 等）
- 样式使用 rpx 单位，组件样式写在单文件 `.vue` 的 `<style>` 标签中
- TypeScript 接口定义和 API 函数放在同一文件

## 常用命令

```bash
cd miniprogram
npm run dev:h5          # H5 开发
npm run type-check      # 类型检查
npm run build:h5        # 生产构建
```
