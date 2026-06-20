package knowledge

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
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
// sqlDB 可选，传入时从 MySQL 加载知识库，否则从内置种子数据加载
func Init(ctx context.Context, sqlDB *sql.DB, llmCfg LLMConfig) error {
	var initErr error
	initOnce.Do(func() {
		db := chromem.NewDB()
		embedFunc := newEmbeddingFunc(llmCfg)
		var err error
		collection, err = db.CreateCollection("pet-knowledge", nil, embedFunc)
		if err != nil {
			initErr = fmt.Errorf("创建 collection 失败: %w", err)
			return
		}
		if sqlDB != nil {
			if err := LoadFromDB(ctx, sqlDB); err != nil {
				initErr = fmt.Errorf("从数据库加载知识库失败: %w", err)
				return
			}
		} else {
			if err := seed(ctx); err != nil {
				initErr = fmt.Errorf("加载内置知识库失败: %w", err)
				return
			}
			log.Printf("内置知识库加载完成，共 %d 条", len(SeedData))
		}
		llm = NewProvider(llmCfg)
	})
	return initErr
}

// GetLLM 获取 LLM 提供者（供 handler 调用）
func GetLLM() LLMProvider {
	return llm
}

// toFloat32 []float64 → []float32
func toFloat32(src []float64) []float32 {
	dst := make([]float32, len(src))
	for i, v := range src {
		dst[i] = float32(v)
	}
	return dst
}

// newEmbeddingFunc 创建 chromem-go 兼容的 embedding 函数
// 支持：（1）Ollama （2）OpenAI/DeepSeek （3）本地 cosine 降级
func newEmbeddingFunc(cfg LLMConfig) chromem.EmbeddingFunc {
	url := cfg.EmbeddingURL
	useOllama := strings.Contains(url, "localhost") || strings.Contains(url, "127.0.0.1") || strings.Contains(url, "0.0.0.0")
	model := cfg.EmbeddingModel
	if model == "" {
		model = "nomic-embed-text"
	}

	return func(ctx context.Context, text string) ([]float32, error) {
		if url == "" {
			return toFloat32(cosineChunkEmbedding(text)), nil
		}

		var body map[string]interface{}
		if useOllama {
			body = map[string]interface{}{"model": model, "prompt": text}
			} else {
			body = map[string]interface{}{"model": model, "input": text}
		}
		data, _ := json.Marshal(body)
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
		if err != nil {
			log.Printf("Embedding 请求创建失败 (%v)，使用本地 embedding", err)
			return toFloat32(cosineChunkEmbedding(text)), nil
		}
		req.Header.Set("Content-Type", "application/json")
		if !useOllama && cfg.APIKey != "" {
			req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Embedding API 调用失败 (%v)，使用本地 embedding", err)
			return toFloat32(cosineChunkEmbedding(text)), nil
		}
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		if len(respBody) == 0 {
			log.Printf("Embedding API 返回空响应，使用本地 embedding")
			return toFloat32(cosineChunkEmbedding(text)), nil
		}

		// Ollama 格式: {"embedding": [0.1, ...]}
		if useOllama {
			var ollamaResp struct {
				Embedding []float64 `json:"embedding"`
			}
			if err := json.Unmarshal(respBody, &ollamaResp); err != nil || len(ollamaResp.Embedding) == 0 {
				log.Printf("Ollama embedding 响应异常，使用本地 embedding")
				return toFloat32(cosineChunkEmbedding(text)), nil
			}
			return toFloat32(ollamaResp.Embedding), nil
		}

		// OpenAI 格式: {"data": [{"embedding": [0.1, ...]}]}
		var openAIResp struct {
			Data []struct {
				Embedding []float64 `json:"embedding"`
			} `json:"data"`
		}
		if err := json.Unmarshal(respBody, &openAIResp); err != nil || len(openAIResp.Data) == 0 {
			log.Printf("Embedding API 响应异常，使用本地 embedding")
			return toFloat32(cosineChunkEmbedding(text)), nil
		}
		return toFloat32(openAIResp.Data[0].Embedding), nil
	}
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

// SearchByBreed 按品种搜索（向量+关键词混合）
func SearchByBreed(ctx context.Context, breed, query string, limit int) ([]Result, error) {
	if limit <= 0 {
		limit = 5
	}
	fullQuery := breed
	if query != "" {
		fullQuery = breed + " " + query
	}
	// 向量搜索
	results, err := Search(ctx, fullQuery, limit*2)
	if err != nil {
		return nil, err
	}
	// 关键词过滤：品种或标签包含搜索词
	breedLow := strings.ToLower(breed)
	var matched []Result
	for _, r := range results {
		if keywordMatch(r, breedLow) || keywordMatch(r, strings.ToLower(query)) {
			matched = append(matched, r)
		}
	}
	// 如果有关键词匹配的结果，用它们
	if len(matched) > 0 {
		if len(matched) > limit {
			matched = matched[:limit]
		}
		return matched, nil
	}
	// 无关键词匹配 → 返回空（不返回模糊结果）
	return nil, nil
}

// keywordMatch 检查结果是否在标题/品种/标签中包含关键字
func keywordMatch(r Result, keyword string) bool {
	if keyword == "" {
		return true
	}
	if strings.Contains(strings.ToLower(r.Title), keyword) {
		return true
	}
	if strings.Contains(strings.ToLower(r.Content), keyword) {
		return true
	}
	for _, t := range r.Tags {
		if strings.Contains(strings.ToLower(t), keyword) {
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

// previewStr 截断字节数组为可读字符串（最长 200 字符）
func previewStr(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	s := string(b)
	if len(s) > 200 {
		s = s[:200]
	}
	return s
}
