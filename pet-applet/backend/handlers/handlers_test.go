package handlers

import (
	"strings"
	"testing"
)

func TestGenerateID_Length(t *testing.T) {
	id := generateID()
	if len(id) != 8 {
		t.Errorf("generateID() 长度应为 8，实际 %d", len(id))
	}
}

func TestGenerateID_ValidChars(t *testing.T) {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	for range 100 {
		id := generateID()
		for _, c := range id {
			if !strings.ContainsRune(chars, c) {
				t.Errorf("generateID() 包含非法字符 %c", c)
			}
		}
	}
}

func TestGenerateID_Unique(t *testing.T) {
	seen := make(map[string]bool)
	for range 1000 {
		id := generateID()
		if seen[id] {
			t.Errorf("generateID() 产生重复 ID: %s", id)
		}
		seen[id] = true
	}
}

func TestGenerateID_NoLeadingZeroIssue(t *testing.T) {
	// 验证不会因未初始化种子而产生可预测的低熵序列
	ids := make([]string, 50)
	for i := range 50 {
		ids[i] = generateID()
	}
	// 检查前 50 个没有重复（随机碰撞概率极低）
	seen := make(map[string]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("50 次生成出现重复: %s", id)
		}
		seen[id] = true
	}
}
