export const meta = {
  name: 'audit',
  description: '全量审计 — 代码质量、测试覆盖、安全扫描三合一',
  phases: [
    { title: '编译检查' },
    { title: '测试覆盖' },
    { title: '审计报告' },
  ],
}

phase('编译检查')
log('🔧 开始编译检查...')

const buildResults = await parallel([
  () => agent('执行后端编译检查并修复：\ncd backend && go vet ./... && go build -o /dev/null ./...\n如果失败，列出错误并修复。', {
    label: 'go: 编译检查',
    phase: '编译检查',
  }),
  () => agent('执行前端编译检查并修复：\ncd miniprogram && npx vue-tsc --noEmit\n如果失败，列出类型错误并修复。', {
    label: 'ts: 类型检查',
    phase: '编译检查',
  }),
  () => agent('执行前端构建检查并修复：\ncd miniprogram && npx vite build --mode production 2>&1\n如果失败，列出构建错误并修复。', {
    label: 'vite: 构建检查',
    phase: '编译检查',
  }),
])

const allBuildOk = buildResults.every(Boolean)
if (!allBuildOk) {
  log('⚠️ 部分编译检查未通过，继续审计但标记风险')
}

phase('测试覆盖')
log('🧪 运行全量测试...')

const testResults = await parallel([
  () => agent('运行后端测试并报告覆盖率：\ncd backend && go test -cover ./...\n解析结果，列出包名、测试通过数、覆盖率百分比。如果有失败测试请修复。', {
    label: 'go: 测试',
    phase: '测试覆盖',
  }),
  () => agent('运行前端测试：\ncd miniprogram && npx vitest run --reporter=verbose 2>&1\n解析结果，列出测试文件数、用例数、通过率。如果有失败测试请修复。', {
    label: 'vitest: 测试',
    phase: '测试覆盖',
  }),
])

phase('审计报告')
log('📊 生成审计报告...')

const report = {
  timestamp: new Date().toISOString(),
  build: {
    backend: buildResults[0] ? '✅ 通过' : '❌ 失败',
    typescript: buildResults[1] ? '✅ 通过' : '❌ 失败',
    frontendBuild: buildResults[2] ? '✅ 通过' : '❌ 失败',
  },
  tests: {
    backend: testResults[0] ? '✅ 通过' : '❌ 失败',
    frontend: testResults[1] ? '✅ 通过' : '❌ 失败',
  },
  overall: allBuildOk && testResults.every(Boolean) ? '✅ 全部通过' : '⚠️ 存在问题',
}

log('========================================')
log('📋 审计报告')
log('========================================')
log('')
log('🔧 编译检查:')
log('   后端编译:     ' + report.build.backend)
log('   类型检查:     ' + report.build.typescript)
log('   前端构建:     ' + report.build.frontendBuild)
log('')
log('🧪 测试覆盖:')
log('   后端测试:     ' + report.tests.backend)
log('   前端测试:     ' + report.tests.frontend)
log('')
log('🏁 总体状态:     ' + report.overall)
log('========================================')

return report
