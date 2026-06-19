export const meta = {
  name: 'parallel-task',
  description: '并行多 agent 执行 — git worktree 隔离 + 自动合并',
  phases: [
    { title: '分析拆解' },
    { title: '环境准备' },
    { title: '并行执行' },
    { title: '集成合并' },
    { title: '代码审查' },
    { title: '清理完成' },
  ],
}

const SUBTASK_SCHEMA = {
  type: 'object',
  properties: {
    subtasks: {
      type: 'array',
      items: {
        type: 'object',
        properties: {
          name: { type: 'string', pattern: '^[a-zA-Z0-9_-]+$', description: '仅允许字母数字下划线连字符' },
          agent: { type: 'string', enum: ['db-dev', 'go-dev', 'uni-app-dev'] },
          target: { type: 'string', description: '修改目标描述' },
          dependsOn: { type: 'array', items: { type: 'string' }, description: '依赖的子任务名列表' },
        },
        required: ['name', 'agent', 'target'],
      },
    },
    apiContract: { type: 'string', description: 'API 契约描述（请求/响应格式）' },
    databaseSchema: { type: 'string', description: '数据库变更 DDL' },
  },
  required: ['subtasks'],
}

phase('分析拆解')
log('🔄 开始分析任务，拆解为可并行子任务...')

const plan = await agent(args?.task || '请描述要执行的任务', {
  label: 'plan-dev: 方案设计',
  phase: '分析拆解',
  schema: SUBTASK_SCHEMA,
})

if (!plan || !plan.subtasks?.length) {
  log('❌ 无法拆解子任务，终止')
  return { success: false, reason: '任务拆解失败' }
}

log('📋 拆解为 ' + plan.subtasks.length + ' 个子任务:')
plan.subtasks.forEach(s => log('   - ' + s.name + ' (' + s.agent + ')'))

phase('环境准备')
log('🔧 准备 git worktree 环境...')

const FEATURE_NAME = args?.feature || 'parallel-' + Date.now()
const BASEDIR = process.cwd()
const WORKTREE_ROOT = BASEDIR + '/../pet-applet-worktrees'

// 创建 worktree 和分支
const setupSteps = plan.subtasks.map(subtask => async () => {
  const branch = 'feat/' + FEATURE_NAME + '/' + subtask.name
  const wtDir = WORKTREE_ROOT + '/' + subtask.name

  await agent(
    '创建 git worktree 用于子任务 "' + subtask.name + '"：\n\n' +
    '1. 如果分支 ' + branch + ' 已存在，删除它\n' +
    '2. 基于 origin/main 创建新分支: git branch ' + branch + ' origin/main\n' +
    '3. 创建 worktree: git worktree add ' + wtDir + ' ' + branch + '\n' +
    '4. 验证: ls ' + wtDir + '/',
    { label: 'setup:' + subtask.name, phase: '环境准备' }
  )
})

await parallel(setupSteps)
log('✅ 所有 worktree 就绪')

phase('并行执行')
log('⚡ 按依赖关系并行执行子任务...')

// 按依赖分组
const executed = new Set()
const results = []

while (executed.size < plan.subtasks.length) {
  const ready = plan.subtasks.filter(s =>
    !executed.has(s.name) &&
    (!s.dependsOn || s.dependsOn.every(d => executed.has(d)))
  )

  if (!ready.length) {
    log('❌ 存在无法满足的依赖关系，终止')
    break
  }

  const batch = ready.map(subtask => async () => {
    log('▶️ 启动 ' + subtask.name + ' (' + subtask.agent + ')')
    const result = await agent(
      '你正在执行子任务 "' + subtask.name + '"。\n' +
      '目标: ' + subtask.target + '\n' +
      'API 契约: ' + (plan.apiContract || '无') + '\n' +
      '数据库变更: ' + (plan.databaseSchema || '无') + '\n' +
      '依赖: ' + ((subtask.dependsOn || []).join(', ') || '无') + '\n\n' +
      '请进入 worktree ' + WORKTREE_ROOT + '/' + subtask.name + '/ 目录执行该子任务。\n' +
      '完成后提交变更到当前分支。',
      {
        label: subtask.agent + ':' + subtask.name,
        phase: '并行执行',
        isolation: 'worktree',
      }
    )

    executed.add(subtask.name)
    return { name: subtask.name, result }
  })

  const batchResults = await parallel(batch.map(f => f))
  results.push(...batchResults.filter(Boolean))
  log('✅ 完成 ' + batch.length + ' 个子任务 (' + executed.size + '/' + plan.subtasks.length + ')')
}

phase('集成合并')
log('🔀 开始合并子任务分支到集成分支...')

const mergeLines = plan.subtasks.map(s => '   - feat/' + FEATURE_NAME + '/' + s.name).join('\n')

const mergeResult = await agent(
  '将各子任务分支合并到集成分支：\n\n' +
  '1. 创建集成分支: git checkout -b feat/integration/' + FEATURE_NAME + ' origin/main\n' +
  '2. 按依赖顺序合并子任务分支:\n' + mergeLines + '\n' +
  '3. 处理冲突（如有）\n' +
  '4. 全量测试:\n' +
  '   - cd ' + BASEDIR + '/backend && go vet ./... && go test ./...\n' +
  '   - cd ' + BASEDIR + '/miniprogram && npx vue-tsc --noEmit\n' +
  '5. 合并到 main: git checkout main && git merge feat/integration/' + FEATURE_NAME + ' --no-ff',
  { label: 'merge-manager: 集成合并', phase: '集成合并' }
)

if (!mergeResult) {
  log('❌ 合并失败，可执行以下命令手动清理：')
  log('   rm -rf ' + WORKTREE_ROOT + '/*')
  log('   git branch -D $(git branch --list \'feat/' + FEATURE_NAME + '/*\')')
  log('   git worktree prune')
  log('   git checkout main')
  return { success: false, reason: '合并冲突或测试失败', feature: FEATURE_NAME }
}

phase('代码审查')
log('🔍 审查合并后的总变更...')

const reviewResult = await agent(
  '审查当前合并后的全部变更（git diff main~1..main）：\n\n' +
  '- 正确性：逻辑是否有缺陷\n' +
  '- 安全性：是否有注入、泄露风险\n' +
  '- 质量：代码是否规范、可维护\n' +
  '- 联动：前后端 API 是否匹配\n' +
  '- 数据库：模型与查询是否同步\n\n' +
  '输出问题清单，按严重程度分级。\n' +
  '严重问题请标注"严重（必须修复）"，建议类请标注"建议（值得改进）"',
  { label: 'code-reviewer: 审查合并', phase: '代码审查' }
)

if (reviewResult && reviewResult.includes('严重（必须修复）')) {
  log('⚠️ 审查发现严重问题，需要修复')
} else {
  log('✅ 审查通过')
}

phase('清理完成')
log('🧹 清理 worktree 和临时分支...')

const wtRemoveLines = plan.subtasks.map(s => '   - git worktree remove ' + WORKTREE_ROOT + '/' + s.name).join('\n')
const branchDeleteLines = plan.subtasks.map(s => '   - git branch -D feat/' + FEATURE_NAME + '/' + s.name).join('\n')

await agent(
  '清理并行任务的临时资源：\n\n' +
  '1. 删除 worktree:\n' + wtRemoveLines + '\n' +
  '2. 删除临时分支:\n' + branchDeleteLines + '\n' +
  '   - git branch -D feat/integration/' + FEATURE_NAME + '\n' +
  '3. git worktree prune',
  { label: 'merge-manager: 清理', phase: '清理完成' }
)

log('✅ 并行任务完成！')
return { success: true, feature: FEATURE_NAME, subtasks: plan.subtasks.length }
