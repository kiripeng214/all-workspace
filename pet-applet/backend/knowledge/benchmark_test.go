package knowledge

import (
	"strings"
	"testing"
)

func BenchmarkCosineChunkEmbedding_Short(b *testing.B) {
	text := "金毛掉毛怎么办"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cosineChunkEmbedding(text)
	}
}

func BenchmarkCosineChunkEmbedding_Long(b *testing.B) {
	text := strings.Repeat("金毛是双层被毛犬种掉毛量大尤其在春秋换毛季。", 50)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cosineChunkEmbedding(text)
	}
}

func BenchmarkKeywordMatch(b *testing.B) {
	r := Result{
		Title:   "金毛掉毛怎么办",
		Content: "金毛是双层被毛犬种，掉毛量大，尤其在春秋换毛季。应对策略：每天梳毛一次（使用针梳+底绒梳），每周洗澡一次辅助去死毛，饮食中添加鱼油有助于皮肤健康。掉毛属正常生理现象，无法完全避免。",
		Tags:    []string{"金毛", "掉毛", "护理"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		keywordMatch(r, "金毛")
		keywordMatch(r, "柯基")
	}
}

func BenchmarkBuildPrompt(b *testing.B) {
	results := make([]Result, 5)
	for i := range results {
		results[i] = Result{
			Title:   "知识标题",
			Content: strings.Repeat("知识内容", 100),
			Score:   0.9,
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buildPrompt("测试问题", results)
	}
}

func BenchmarkToFloat32(b *testing.B) {
	src := make([]float64, 768)
	for i := range src {
		src[i] = float64(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		toFloat32(src)
	}
}
