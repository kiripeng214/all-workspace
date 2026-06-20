package knowledge

// KnowledgeEntry 知识条目
type KnowledgeEntry struct {
	Breed   string   // 适用品种（空=通用）
	Title   string   // 标题
	Content string   // 正文
	Tags    []string // 标签
}

// SeedData 内置种子知识库（降级用，正常由 MySQL 提供）
// 数据已迁移到 MySQL knowledge_entries 表 + backend/data/seed_knowledge.json
var SeedData []KnowledgeEntry
