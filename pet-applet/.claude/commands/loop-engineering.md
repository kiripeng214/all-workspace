---
name: project:loop-engineering
description: 循环工程 — 自动审查并调整 Loop Engine 体系
---

# /project:loop-engineering

对项目的 Loop Engine（`.claude/` 配置体系）进行全面审查，自动检测结构缺失、配置漂移、关联断裂等问题，并执行修复调整。

## 审查维度

### 1. 结构完整性
检查所有目录和文件是否存在：

| 检查项 | 命令 |
|--------|------|
| commands | `ls .claude/commands/*.md` |
| agents | `ls .claude/agents/*.md` |
| rules | `find .claude/rules -name "RULE.md"` |
| skills | `find .claude/skills -name "SKILL.md"` |
| workflows | `ls .claude/workflows/*.js 2>/dev/null` |

### 2. 配置漂移检测
检查以下一致性：

- **Agent 引用一致性** — 检查 `workflows/*.js` 中引用的 agent 名称是否在 `.claude/agents/` 中有定义文件
- **Command frontmatter** — 检查每个 `.md` 文件是否有 `name` / `description` / `trigger` 字段，且 mapping 正确
- **Agent tools 校验** — 检查 agent 声明的 tools 是否为有效工具名
- **Skill 引用** — 检查 settings.json 的 `skills` 路径映射是否与 skills 目录一致
- **Rule 目录** — 检查 rules 路径映射是否与 rules 目录一致

### 3. 技能/命令/Agent 关联分析
检查以下关联是否完整：

- 每个 `/project:*` 命令是否有对应的 skill 或 agent 定义
- 每个 skill 是否有对应的命令入口（方便 `/` 快捷调用）
- pre-commit trigger 的 skill 是否有对应的 hook 配置

### 4. 项目适配性检查
检查 Loop Engine 配置是否与实际项目匹配：

- 前端框架检测（uni-app / Vue 3 / Vite）
- 后端框架检测（Go / Gin）
- 测试框架检测（vitest / go test）
- 构建命令检测（`npm run dev:h5` / `go build`）

## 执行流程

### 阶段 1：快速扫描

```bash
# 结构清单
echo "=== Commands ===" && ls .claude/commands/*.md
echo "=== Agents ===" && ls .claude/agents/*.md
echo "=== Rules ===" && find .claude/rules -name "RULE.md"
echo "=== Skills ===" && find .claude/skills -name "SKILL.md"
echo "=== Workflows ===" && ls .claude/workflows/*.js 2>/dev/null
```

### 阶段 2：深度审查

使用 `loop-auditor` agent 对发现的缺口进行逐项分析，对每个问题给出：

```
缺失类型: [结构缺失 | 配置漂移 | 引用断裂]
问题详情: ...
建议修复: ...
优先级:   [高 | 中 | 低]
```

### 阶段 3：执行修复

对可自动修复的问题直接执行：

| 问题类型 | 自动修复方式 |
|----------|-------------|
| 缺失命令文件 | 根据 skill 内容自动生成命令文件 |
| 缺失 frontmatter | 自动补全 name / description 字段 |
| Agent tools 缺失 | 补充合理的默认 tools 列表 |
| Workflow 引用断裂 | 补充缺失的 agent 定义文件 |
| 路径映射不一致 | 更新 settings.json |

## 输出格式

```
========================================
🔧 Loop Engineering 报告
========================================

📂 结构状态:  ✅ 完整 | ⚠️ 缺失 N 项
🔄 配置漂移:  ✅ 一致 | ⚠️ 漂移 N 处
🔗 关联断裂:  ✅ 完整 | ⚠️ 断裂 N 处
🎯 项目匹配:  ✅ 适配 | ⚠️ 不匹配 N 处

📋 待修复项（N 项）:
  [高] ...（自动修复 ✓ | 需手动 ✗）
  [中] ...

✨ 优化建议:
  - ...

========================================
```

## 使用方式

```bash
# 审查并自动修复所有问题
/project:loop-engineering

# 仅审查不做修改（review-only 模式待实现）
```
