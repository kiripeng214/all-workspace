---
name: project:frontend-dev
description: 启动前端 H5 开发服务器
---

# /project:frontend-dev

启动 uni-app H5 开发模式（Vite 热更新）。

## 执行步骤

1. 进入 `miniprogram/` 目录
2. 执行 `npm run dev:h5`

## 其他平台

若需要其他平台，替换为对应命令：

| 平台 | 命令 |
|---|---|
| 微信小程序 | `npm run dev:mp-weixin` |
| 支付宝小程序 | `npm run dev:mp-alipay` |
| 百度小程序 | `npm run dev:mp-baidu` |
| 抖音小程序 | `npm run dev:mp-toutiao` |
| QQ 小程序 | `npm run dev:mp-qq` |
| 快应用 | `npm run dev:quickapp-webview` |
| Harmony 鸿蒙 | `npm run dev:mp-harmony` |

## 验证

终端输出 Vite dev server URL（默认 `http://localhost:5173`），浏览器打开可见宠物列表页。
