package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pet-applet-backend/knowledge"
)

// SearchKnowledge 搜索知识库
func SearchKnowledge(c *gin.Context) {
	query := c.Query("q")
	breed := c.Query("breed")

	if query == "" && breed == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入搜索关键词或选择品种"})
		return
	}

	var results []knowledge.Result
	var err error

	if breed != "" {
		results, err = knowledge.SearchByBreed(c.Request.Context(), breed, query, 5)
	} else {
		results, err = knowledge.Search(c.Request.Context(), query, 5)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 调用 LLM 生成回答
	llmResp, err := knowledge.QueryLLM(c.Request.Context(), query, results)
	if err != nil {
		// LLM 失败返回纯检索结果
		c.JSON(http.StatusOK, gin.H{
			"results": results,
			"answer":  nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"answer":  llmResp,
	})
}
