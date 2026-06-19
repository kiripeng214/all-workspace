---
name: project:deploy-build
description: 生产构建 — 后端 Docker 镜像 + 前端微信小程序打包
---

# /project:deploy-build

一键构建生产产物：

| 产物 | 方式 | 输出 |
|------|------|------|
| 后端镜像 | Docker 多阶段构建 | `pet-applet-backend:latest` |
| 前端小程序 | `uni build -p mp-weixin` | `miniprogram/dist/build/mp-weixin/` |

## 执行步骤

### 1. 构建后端 Docker 镜像

```bash
cd backend
docker build -t pet-applet-backend:latest .
```

验证：`docker images pet-applet-backend`

### 2. 打包前端微信小程序

```bash
cd miniprogram
npm run build:mp-weixin
```

产物路径：`miniprogram/dist/build/mp-weixin/`

### 3. 验证产物

```bash
# 检查镜像
docker images pet-applet-backend:latest

# 检查小程序产物
ls -la miniprogram/dist/build/mp-weixin/
```

## 注意

- 后端 Dockerfile 位置：`backend/Dockerfile`（多阶段构建，alpine 运行）
- 前端需在 `manifest.json` 中正确配置微信小程序的 AppID
- 首次使用需要先登录微信开发者工具关联项目
