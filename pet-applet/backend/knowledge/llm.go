package knowledge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// LLMProvider LLM 提供者接口
type LLMProvider interface {
	Ask(ctx context.Context, prompt string) (string, error)
}

// LLMConfig LLM 配置
type LLMConfig struct {
	Provider string
	APIKey   string
	APIURL   string
	Model    string
}

// LLMResponse LLM 回答
type LLMResponse struct {
	Answer  string   `json:"answer"`
	Sources []string `json:"sources"`
}

// OpenAIProvider OpenAI 兼容 API 实现
type OpenAIProvider struct {
	apiKey string
	apiURL string
	model  string
}

func (p *OpenAIProvider) Ask(ctx context.Context, prompt string) (string, error) {
	body := map[string]interface{}{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 1024,
	}
	data, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", p.apiURL, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil || len(result.Choices) == 0 {
		return "", fmt.Errorf("OpenAI 响应解析失败: %s", string(respBody))
	}
	return result.Choices[0].Message.Content, nil
}

// NewProvider 根据配置创建 LLM 提供者
func NewProvider(cfg LLMConfig) LLMProvider {
	if cfg.APIKey == "" || cfg.APIURL == "" {
		log.Println("LLM 未配置，使用降级模式")
		return nil
	}
	switch cfg.Provider {
	case "anthropic":
		return &AnthropicProvider{apiKey: cfg.APIKey, apiURL: cfg.APIURL, model: cfg.Model}
	default:
		return &OpenAIProvider{apiKey: cfg.APIKey, apiURL: cfg.APIURL, model: cfg.Model}
	}
}

// buildPrompt 构建 RAG 提示词
func buildPrompt(query string, results []Result) (string, []string) {
	var contextParts []string
	var sources []string
	seen := make(map[string]bool)

	for _, r := range results {
		if r.Score < 0.3 {
			continue
		}
		contextParts = append(contextParts, fmt.Sprintf("## %s\n%s", r.Title, r.Content))
		if !seen[r.Title] {
			sources = append(sources, r.Title)
			seen[r.Title] = true
		}
	}

	contextStr := strings.Join(contextParts, "\n\n")
	if contextStr == "" {
		return "", sources
	}

	prompt := fmt.Sprintf(`你是一个宠物养护专家。基于以下知识库内容回答用户的问题。

## 知识库内容
%s

## 用户问题
%s

请用中文回答，语言通俗易懂，不要编造知识库中没有的信息。如果知识库内容不足以回答问题，请如实告知。`, contextStr, query)

	return prompt, sources
}

// QueryLLM 使用 LLM provider 生成回答
func QueryLLM(ctx context.Context, query string, results []Result, provider LLMProvider) *LLMResponse {
	prompt, sources := buildPrompt(query, results)
	if prompt == "" {
		return fallbackResponse(query, results)
	}

	if provider == nil {
		return fallbackResponse(query, results)
	}

	answer, err := provider.Ask(ctx, prompt)
	if err != nil {
		log.Printf("LLM 调用失败: %v，降级为检索结果", err)
		return fallbackResponse(query, results)
	}

	return &LLMResponse{
		Answer:  answer,
		Sources: sources,
	}
}

// fallbackResponse 无 LLM 时返回纯检索结果
func fallbackResponse(query string, results []Result) *LLMResponse {
	if len(results) == 0 {
		return &LLMResponse{
			Answer:  "抱歉，未找到相关知识。请尝试换个关键词搜索。",
			Sources: nil,
		}
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("以下是与「%s」相关的知识：\n\n", query))
	var sources []string
	seen := make(map[string]bool)

	for i, r := range results {
		if i >= 3 {
			break
		}
		b.WriteString(fmt.Sprintf("**%s**\n%s\n\n", r.Title, r.Content))
		if !seen[r.Title] {
			sources = append(sources, r.Title)
			seen[r.Title] = true
		}
	}

	return &LLMResponse{
		Answer:  b.String(),
		Sources: sources,
	}
}
