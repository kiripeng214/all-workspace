---
name: project:backend-start
description: 构建并启动 Go 后端服务
---

# /project:backend-start

构建 backend/ 目录下的 Go 项目并启动服务。

## 执行步骤

1. 进入 `backend/` 目录
2. 编译：`go build -o pet-applet-server .`
3. 启动服务：`./pet-applet-server`（后台运行）

## 前置条件

- 已安装 Go 1.26+
- MySQL 服务运行中
- `backend/config/config.yaml` 中数据库连接配置正确

## 验证

服务启动后在终端输出 `🐾 服务启动于 http://localhost:3000`，访问 `http://localhost:3000/api/pets` 返回 JSON。
