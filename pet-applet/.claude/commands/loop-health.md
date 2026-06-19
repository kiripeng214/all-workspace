---
name: project:loop-health
description: Loop Engine 健康检查 — 审计 .claude 体系搭建完整度，输出改进建议
---

# /project:loop-health

一键审计当前项目的 Loop Engine 搭建状态。

## 执行流程

### 1. 结构完整性检查

| 检查项 | 命令 |
|--------|------|
| commands 目录 | `ls .claude/commands/*.md` |
| agents 目录 | `ls .claude/agents/*.md` |
| rules 目录 | `find .claude/rules -name "RULE.md"` |
| skills 目录 | `find .claude/skills -name "SKILL.md"` |
| workflows 目录 | `ls .claude/workflows/*.js 2>/dev/null` |
| settings.json | `cat .claude/settings.json` |

### 2. 测试覆盖检查

| 检查项 | 命令 |
|--------|------|
| 前端测试 | `cd miniprogram && npx vitest run 2>&1 \| tail -5` |
| 前端类型检查 | `cd miniprogram && npx vue-tsc --noEmit 2>&1` |
| 后端测试 | `cd backend && go test ./... 2>&1` |
| 后端编译 | `cd backend && go vet ./... 2>&1` |

### 3. 输出审计摘要

- 结构清单（存在的 / 缺失的）
- 测试状态（通过 / 失败）
- 构建状态（通过 / 失败）
- 改进建议列表

> 输出结果后，询问是否要执行改进项。
