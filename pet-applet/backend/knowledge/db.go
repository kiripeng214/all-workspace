package knowledge

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/philippgille/chromem-go"
)

// LoadFromDB 从 MySQL 加载所有知识条目到 chromem-go
func LoadFromDB(ctx context.Context, db *sql.DB) error {
	rows, err := db.Query("SELECT breed, title, content, tags FROM knowledge_entries ORDER BY id")
	if err != nil {
		return fmt.Errorf("查询知识条目失败: %w", err)
	}
	defer rows.Close()

	type entry struct {
		breed string
		title string
		content string
		tags   string
	}

	var entries []entry
	for rows.Next() {
		var e entry
		if err := rows.Scan(&e.breed, &e.title, &e.content, &e.tags); err != nil {
			return fmt.Errorf("扫描知识条目失败: %w", err)
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	if len(entries) == 0 {
		return fmt.Errorf("knowledge_entries 表为空，请先执行迁移")
	}

	for i, e := range entries {
		text := e.title + "\n" + e.content
		id := fmt.Sprintf("pet_knowledge_%04d", i)

		// 解析标签 (逗号分隔)
		tagList := strings.Split(e.tags, ",")
		var cleanTags []string
		for _, t := range tagList {
			t = strings.TrimSpace(t)
			if t != "" {
				cleanTags = append(cleanTags, t)
			}
		}

		metadata := map[string]string{
			"title": e.title,
			"breed": e.breed,
			"tags":  strings.Join(cleanTags, ","),
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

	if err := rows.Err(); err != nil {
		return err
	}

	log.Printf("知识库从 MySQL 加载完成，共 %d 条", len(entries))
	return nil
}
