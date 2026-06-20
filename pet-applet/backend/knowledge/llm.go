package knowledge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// LLMResponse LLM 回答
type LLMResponse struct {
	Answer  string   `json:"answer"`
	Sources []string `json:"sources"`
}

// QueryLLM 调用 LLM API 生成回答
func QueryLLM(ctx context.Context, query string, results []Result) (*LLMResponse, error) {
	apiKey := os.Getenv("LLM_API_KEY")
	apiURL := os.Getenv("LLM_API_URL")
	model := os.Getenv("LLM_MODEL")

	if apiKey == "" || apiURL == "" {
		return fallbackResponse(query, results), nil
	}
	if model == "" {
		model = "claude-sonnet-4-20250514"
	}

	// 构建上下文
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
		return fallbackResponse(query, results), nil
	}

	prompt := fmt.Sprintf(`你是一个宠物养护专家。基于以下知识库内容回答用户的问题。

## 知识库内容
%s

## 用户问题
%s

请用中文回答，语言通俗易懂，不要编造知识库中没有的信息。如果知识库内容不足以回答问题，请如实告知。`, contextStr, query)

	// 调用 OpenAPI 兼容的 LLM API
	body := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": 1024,
	}

	data, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(data))
	if err != nil {
		return fallbackResponse(query, results), nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fallbackResponse(query, results), nil
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var apiResult struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &apiResult); err != nil || len(apiResult.Choices) == 0 {
		// 尝试 Anthropic 格式
		var anthropicResult struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		}
		if err := json.Unmarshal(respBody, &anthropicResult); err == nil && len(anthropicResult.Content) > 0 {
			return &LLMResponse{
				Answer:  anthropicResult.Content[0].Text,
				Sources: sources,
			}, nil
		}
		return fallbackResponse(query, results), nil
	}

	return &LLMResponse{
		Answer:  apiResult.Choices[0].Message.Content,
		Sources: sources,
	}, nil
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
