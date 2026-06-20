package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"

	"pet-applet-backend/config"
)

var DB *sql.DB

func Init(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("数据库无法访问: %v", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	log.Println("数据库连接成功")

	if err := goose.SetDialect("mysql"); err != nil {
		log.Fatalf("设置 goose dialect 失败: %v", err)
	}
	dir := resolveMigrationsDir()
	if err := goose.Up(DB, dir); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
}

// resolveMigrationsDir 按优先级查找迁移目录
func resolveMigrationsDir() string {
	// 1. 环境变量指定
	if d := os.Getenv("MIGRATIONS_DIR"); d != "" {
		return d
	}
	// 2. 可执行文件同级（生产 Docker 镜像）
	if execPath, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(execPath), "migrations")
		if info, statErr := os.Stat(candidate); statErr == nil && info.IsDir() {
			return candidate
		}
	}
	// 3. 向上查找后端项目根目录的 migrations/（开发 / 测试）
	wd, _ := os.Getwd()
	for dir := wd; dir != "/"; dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, "migrations")
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			// 确认是后端项目的 migrations（包含 001_create_pets.sql）
			if _, err := os.Stat(filepath.Join(candidate, "001_create_pets.sql")); err == nil {
				return candidate
			}
		}
	}
	// 4. 最终后备
	return "migrations"
}
