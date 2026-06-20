package knowledge

import (
	"math"
	"strings"
	"testing"
)

func TestToFloat32(t *testing.T) {
	src := []float64{1.0, 2.5, -3.14, 0}
	dst := toFloat32(src)
	if len(dst) != 4 {
		t.Fatalf("len 应为 4，实际 %d", len(dst))
	}
	if dst[0] != float32(1.0) || dst[2] != float32(-3.14) {
		t.Errorf("值转换错误: %v", dst)
	}
}

func TestToFloat32_Empty(t *testing.T) {
	if dst := toFloat32(nil); len(dst) != 0 {
		t.Error("nil 输入应返回空切片")
	}
	if dst := toFloat32([]float64{}); len(dst) != 0 {
		t.Error("空切片应返回空切片")
	}
}

func TestCosineChunkEmbedding_Dim(t *testing.T) {
	vec := cosineChunkEmbedding("金毛掉毛怎么办")
	if len(vec) != 64 {
		t.Errorf("维度应为 64，实际 %d", len(vec))
	}
}

func TestCosineChunkEmbedding_Normalized(t *testing.T) {
	vec := cosineChunkEmbedding("测试文本")
	var sum float64
	for _, v := range vec {
		sum += v * v
	}
	norm := math.Sqrt(sum)
	if math.Abs(norm-1.0) > 1e-6 {
		t.Errorf("未归一化，norm=%f", norm)
	}
}

func TestCosineChunkEmbedding_Deterministic(t *testing.T) {
	v1 := cosineChunkEmbedding("柯基腰椎保护")
	v2 := cosineChunkEmbedding("柯基腰椎保护")
	for i := range v1 {
		if v1[i] != v2[i] {
			t.Errorf("相同输入应产生相同 embedding，位置 %d 不同", i)
			break
		}
	}
}

func TestCosineChunkEmbedding_Different(t *testing.T) {
	v1 := cosineChunkEmbedding("金毛掉毛")
	v2 := cosineChunkEmbedding("猫咪喝水")
	same := true
	for i := range v1 {
		if v1[i] != v2[i] {
			same = false
			break
		}
	}
	if same {
		t.Error("不同输入应产生不同 embedding")
	}
}

func TestKeywordMatch(t *testing.T) {
	r := Result{
		Title:   "金毛掉毛怎么办",
		Content: "金毛是双层被毛犬种",
		Tags:    []string{"金毛", "掉毛"},
	}
	tests := []struct {
		keyword string
		want    bool
	}{
		{"金毛", true},
		{"掉毛", true},
		{"双层", true},
		{"柯基", false},
		{"猫咪", false},
		{"", true},
	}
	for _, tt := range tests {
		got := keywordMatch(r, tt.keyword)
		if got != tt.want {
			t.Errorf("keywordMatch(%q) = %v, want %v", tt.keyword, got, tt.want)
		}
	}
}

func TestKeywordMatch_ContentOnly(t *testing.T) {
	r := Result{
		Title:   "普通标题",
		Content: "这里提到了英短猫的养护知识",
		Tags:    []string{"通用"},
	}
	if !keywordMatch(r, "英短") {
		t.Error("应在 content 中匹配到英短")
	}
}

func TestPreviewStr_Short(t *testing.T) {
	s := previewStr([]byte("hello"))
	if s != "hello" {
		t.Errorf("短文本不应截断: %s", s)
	}
}

func TestPreviewStr_Long(t *testing.T) {
	long := strings.Repeat("a", 500)
	s := previewStr([]byte(long))
	if len(s) != 200 {
		t.Errorf("长文本应截断到 200，实际 %d", len(s))
	}
}

func TestPreviewStr_Empty(t *testing.T) {
	if s := previewStr(nil); s != "" {
		t.Errorf("nil 应返回空字符串，实际 %q", s)
	}
	if s := previewStr([]byte{}); s != "" {
		t.Errorf("空切片应返回空字符串，实际 %q", s)
	}
}

func TestBuildPrompt_WithResults(t *testing.T) {
	results := []Result{
		{Title: "标题1", Content: "内容1", Tags: []string{"tag1"}, Score: 0.9},
		{Title: "标题2", Content: "内容2", Tags: []string{"tag2"}, Score: 0.5},
	}
	prompt, sources := buildPrompt("测试问题", results)
	if prompt == "" {
		t.Fatal("有结果时应生成 prompt")
	}
	if len(sources) != 2 {
		t.Errorf("应有 2 个来源，实际 %d", len(sources))
	}
	if !strings.Contains(prompt, "测试问题") {
		t.Error("prompt 应包含用户问题")
	}
	if !strings.Contains(prompt, "标题1") {
		t.Error("prompt 应包含知识标题")
	}
}

func TestBuildPrompt_LowScore(t *testing.T) {
	results := []Result{
		{Title: "无关", Content: "无关内容", Score: 0.1},
	}
	prompt, sources := buildPrompt("问题", results)
	if prompt != "" {
		t.Error("低分结果应返回空 prompt")
	}
	if sources != nil {
		t.Error("低分结果应返回 nil sources")
	}
}

func TestBuildPrompt_Empty(t *testing.T) {
	prompt, sources := buildPrompt("问题", nil)
	if prompt != "" {
		t.Error("空结果应返回空 prompt")
	}
	if sources != nil {
		t.Error("空结果应返回 nil sources")
	}
}

func TestFallbackResponse_WithResults(t *testing.T) {
	results := []Result{
		{Title: "金毛掉毛", Content: "每天梳毛一次", Score: 0.8},
	}
	resp := fallbackResponse("金毛", results)
	if resp == nil {
		t.Fatal("fallbackResponse 不应返回 nil")
	}
	if !strings.Contains(resp.Answer, "金毛") {
		t.Error("回答应包含查询关键词")
	}
	if len(resp.Sources) != 1 {
		t.Errorf("应有 1 个来源，实际 %d", len(resp.Sources))
	}
}

func TestFallbackResponse_Empty(t *testing.T) {
	resp := fallbackResponse("赤狐", nil)
	if resp == nil {
		t.Fatal("fallbackResponse 不应返回 nil")
	}
	if !strings.Contains(resp.Answer, "未找到") {
		t.Error("空结果应提示未找到")
	}
}

func TestFallbackResponse_LimitResults(t *testing.T) {
	results := make([]Result, 10)
	for i := range results {
		results[i] = Result{Title: "知识", Content: "内容", Score: 0.9}
	}
	resp := fallbackResponse("测试", results)
	// 最多显示 3 条
	count := strings.Count(resp.Answer, "**")
	if count > 6 { // 3 * 2 (**)
		t.Errorf("应最多显示 3 条，实际显示了更多")
	}
}
