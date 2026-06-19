export const meta = {
  name: 'deep-review',
  description: '多维度代码审查 — 正确性、安全、性能、前端联动四维度并行审查',
  phases: [
    { title: '准备' },
    { title: '并行审查' },
    { title: '汇总' },
  ],
}

const FINDING_SCHEMA = {
  type: 'object',
  properties: {
    findings: {
      type: 'array',
      items: {
        type: 'object',
        properties: {
          severity: { type: 'string', enum: ['严重', '一般', '建议'] },
          file: { type: 'string' },
          line: { type: 'string' },
          title: { type: 'string' },
          detail: { type: 'string' },
        },
        required: ['severity', 'file', 'title', 'detail'],
      },
    },
  },
  required: ['findings'],
}

phase('准备')
log('🔍 开始多维度代码审查...')

const DIMENSIONS = [
  {
    key: 'correctness',
    prompt: '审查当前 git diff 的正确性。重点关注：\n- 逻辑缺陷（空指针、边界条件、竞态）\n- API 参数校验缺失\n- 数据库查询条件错误\n- 类型断言风险\n- 数据流是否完整（创建→读取→更新→删除）\n\n列出所有正确性问题。',
  },
  {
    key: 'security',
    prompt: '审查当前 git diff 的安全性。重点关注：\n- SQL 注入风险\n- XSS / 注入风险\n- 敏感信息泄露（日志打印密码/token）\n- 越权访问（缺少权限校验）\n- 用户输入未过滤\n\n列出所有安全问题。',
  },
  {
    key: 'performance',
    prompt: '审查当前 git diff 的性能。重点关注：\n- N+1 查询\n- 不必要的重复计算\n- 大对象/循环中的性能隐患\n- 前端组件重复渲染\n- 冗余网络请求\n\n列出所有性能问题。',
  },
  {
    key: 'frontend-backend',
    prompt: '审查当前 git diff 的前后端联动。重点关注：\n- API 路径与 handler 路由是否匹配\n- 请求/响应字段名是否一致\n- 状态码处理是否完整\n- 新增字段在前后端是否都已处理\n- 错误处理是否在前端有对应 UI 反馈\n\n列出所有联动问题。',
  },
]

phase('并行审查')

const results = await pipeline(
  DIMENSIONS,
  dim => agent(dim.prompt, {
    label: 'review:' + dim.key,
    phase: '并行审查',
    schema: FINDING_SCHEMA,
  }),
)

const allFindings = results
  .filter(Boolean)
  .flatMap((r, i) => (r.findings || []).map(f => ({ ...f, dimension: DIMENSIONS[i].key })))

phase('汇总')

const severityOrder = { '严重': 0, '一般': 1, '建议': 2 }
allFindings.sort((a, b) => (severityOrder[a.severity] ?? 9) - (severityOrder[b.severity] ?? 9))

const serious = allFindings.filter(f => f.severity === '严重')
const normal = allFindings.filter(f => f.severity === '一般')
const suggestion = allFindings.filter(f => f.severity === '建议')

log('📋 审查结果汇总:')
log('   严重: ' + serious.length + ' | 一般: ' + normal.length + ' | 建议: ' + suggestion.length + ' | 总计: ' + allFindings.length)

serious.forEach(f => log('   🔴 [' + f.dimension + '] ' + f.title + ' — ' + f.file))
normal.forEach(f => log('   🟡 [' + f.dimension + '] ' + f.title + ' — ' + f.file))
suggestion.forEach(f => log('   🔵 [' + f.dimension + '] ' + f.title + ' — ' + f.file))

return {
  total: allFindings.length,
  serious: serious.length,
  normal: normal.length,
  suggestion: suggestion.length,
  findings: allFindings,
}
