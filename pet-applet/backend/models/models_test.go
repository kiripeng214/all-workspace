package models

import (
	"encoding/json"
	"testing"
)

func TestPet_JSONRoundTrip(t *testing.T) {
	p := Pet{
		ID:        "test1234",
		Avatar:    "🐶",
		Name:      "小狗",
		Breed:     "金毛",
		Birthday:  "2024-01-01",
		Weight:    "15kg",
		Notes:     "很乖",
		CreatedAt: 1700000000000,
	}

	data, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("JSON 序列化失败: %v", err)
	}

	var p2 Pet
	if err := json.Unmarshal(data, &p2); err != nil {
		t.Fatalf("JSON 反序列化失败: %v", err)
	}

	if p2.ID != p.ID || p2.Name != p.Name {
		t.Errorf("JSON round-trip 失败: got %+v, want %+v", p2, p)
	}
}

func TestFeedingSchedule_JSONRoundTrip(t *testing.T) {
	s := FeedingSchedule{
		ID:       "sched001",
		PetID:    "pet001",
		Time:     "08:00",
		FoodType: "粮食",
		Amount:   "一份",
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("JSON 序列化失败: %v", err)
	}

	var s2 FeedingSchedule
	if err := json.Unmarshal(data, &s2); err != nil {
		t.Fatalf("JSON 反序列化失败: %v", err)
	}

	if s2.ID != s.ID || s2.Time != s.Time {
		t.Errorf("JSON round-trip 失败: got %+v, want %+v", s2, s)
	}
}

func TestFeedingRecord_DefaultValues(t *testing.T) {
	// 验证空字段时间 JSON 序列化
	r := FeedingRecord{
		ID:    "rec001",
		PetID: "pet001",
		Time:  "12:00",
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Fatalf("JSON 序列化失败: %v", err)
	}

	var r2 FeedingRecord
	if err := json.Unmarshal(data, &r2); err != nil {
		t.Fatalf("JSON 反序列化失败: %v", err)
	}

	if r2.Notes != "" {
		t.Errorf("空字段 Notes 应序列化为空字符串，实际: %q", r2.Notes)
	}
}
