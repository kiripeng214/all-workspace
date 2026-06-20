package migration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func Run(db *sql.DB, dir string) error {
	if err := ensureTrackTable(db); err != nil {
		return fmt.Errorf("创建迁移追踪表失败: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return fmt.Errorf("读取迁移目录失败: %w", err)
	}
	if len(files) == 0 {
		log.Println("未发现迁移文件")
		return nil
	}

	sort.Strings(files)

	done, err := getExecuted(db)
	if err != nil {
		return err
	}

	for _, f := range files {
		name := filepath.Base(f)
		if done[name] {
			continue
		}

		sqlBytes, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("读取 %s 失败: %w", name, err)
		}

		statements := strings.TrimSpace(string(sqlBytes))
		if statements == "" {
			continue
		}

		if _, err := db.Exec(statements); err != nil {
			return fmt.Errorf("执行迁移 %s 失败: %w", name, err)
		}

		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", name); err != nil {
			return fmt.Errorf("记录迁移 %s 失败: %w", name, err)
		}

		log.Printf("迁移 %s 完成", name)
	}

	return nil
}

func ensureTrackTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(100) PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`)
	return err
}

func getExecuted(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, fmt.Errorf("查询迁移记录失败: %w", err)
	}
	defer rows.Close()

	done := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		done[v] = true
	}
	return done, rows.Err()
}
