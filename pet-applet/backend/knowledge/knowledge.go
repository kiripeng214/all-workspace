package knowledge

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/philippgille/chromem-go"
)

var (
	collection *chromem.Collection
	llm        LLMProvider
	initOnce   sync.Once
)

// Result 检索结果
type Result struct {
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Tags    []string `json:"tags"`
	Score   float64 `json:"score"`
}

// Init 初始化向量知识库和 LLM 提供者
func Init(ctx context.Context, llmCfg LLMConfig) error {
	var initErr error
	initOnce.Do(func() {
		db := chromem.NewDB()
		var err error
		collection, err = db.CreateCollection("pet-knowledge", nil, nil)
		if err != nil {
			initErr = fmt.Errorf("创建 collection 失败: %w", err)
			return
		}
		if err := seed(ctx); err != nil {
			initErr = fmt.Errorf("初始化知识库失败: %w", err)
			return
		}
		llm = NewProvider(llmCfg)
		log.Printf("知识库加载完成，共 %d 条", len(SeedData))
	})
	return initErr
}

// GetLLM 获取 LLM 提供者（供 handler 调用）
func GetLLM() LLMProvider {
	return llm
}

func seed(ctx context.Context) error {
	for i, entry := range SeedData {
		text := entry.Title + "\n" + entry.Content
		id := fmt.Sprintf("pet_knowledge_%04d", i)
		metadata := map[string]string{
			"title": entry.Title,
			"breed": entry.Breed,
			"tags":  strings.Join(entry.Tags, ","),
		}
		doc := chromem.Document{
			ID:       id,
			Metadata: metadata,
			Content:  text,
		}
		if err := collection.AddDocument(ctx, doc); err != nil {
			return fmt.Errorf("添加文档 %s 失败: %w", id, err)
		}
	}
	return nil
}

// Search 搜索知识库
func Search(ctx context.Context, query string, limit int) ([]Result, error) {
	if collection == nil {
		return nil, fmt.Errorf("知识库未初始化")
	}
	if limit <= 0 {
		limit = 5
	}

	results, err := collection.Query(ctx, query, limit, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %w", err)
	}

	var out []Result
	for _, r := range results {
		tags := strings.Split(r.Metadata["tags"], ",")
		// 过滤空标签
		var clean []string
		for _, t := range tags {
			if t != "" {
				clean = append(clean, t)
			}
		}
		out = append(out, Result{
			Title:   r.Metadata["title"],
			Content: r.Content,
			Tags:    clean,
			Score:   float64(r.Similarity),
		})
	}
	return out, nil
}

// SearchByBreed 按品种搜索
func SearchByBreed(ctx context.Context, breed, query string, limit int) ([]Result, error) {
	fullQuery := breed
	if query != "" {
		fullQuery = breed + " " + query
	}
	results, err := Search(ctx, fullQuery, limit)
	if err != nil {
		return nil, err
	}
	// 优先返回匹配品种的结果
	sortByBreed(results, breed)
	return results, nil
}

func sortByBreed(results []Result, breed string) {
	// 简单冒泡：品种匹配的排前面
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			ai := containsTag(results[i].Tags, breed)
			aj := containsTag(results[j].Tags, breed)
			if !ai && aj {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}

func containsTag(tags []string, target string) bool {
	target = strings.ToLower(target)
	for _, t := range tags {
		if strings.Contains(strings.ToLower(t), target) {
			return true
		}
	}
	return false
}

// GetCollection 获取集合（供 handler 使用）
func GetCollection() *chromem.Collection {
	return collection
}

// cosineChunkEmbedding 简单的分块文本嵌入（纯 Go 实现，无需外部 API）
func cosineChunkEmbedding(text string) []float64 {
	// 将文本分成字符块，统计每个 unicode 块的频率
	dim := 64
	vec := make([]float64, dim)
	text = strings.ToLower(text)

	for _, r := range text {
		idx := int(r) % dim
		vec[idx]++
	}

	// L2 归一化
	var sum float64
	for _, v := range vec {
		sum += v * v
	}
	norm := math.Sqrt(sum)
	if norm > 0 {
		for i := range vec {
			vec[i] /= norm
		}
	}
	return vec
}

// EmbedContent 使用本地嵌入（兼容 chromem-go）
func EmbedContent(text string) ([]float64, error) {
	return cosineChunkEmbedding(text), nil
}

// EstimateTokens 粗略估算 token 数
func EstimateTokens(text string) int {
	return utf8.RuneCountInString(text) / 2
}
