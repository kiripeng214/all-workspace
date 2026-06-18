package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

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
	migrate()
}

func migrate() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS pets (
			id VARCHAR(36) PRIMARY KEY,
			avatar VARCHAR(10) NOT NULL DEFAULT '🐾',
			name VARCHAR(100) NOT NULL,
			breed VARCHAR(100) DEFAULT '',
			birthday VARCHAR(20) DEFAULT '',
			weight VARCHAR(20) DEFAULT '',
			notes TEXT,
			created_at BIGINT NOT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS feeding_schedules (
			id VARCHAR(36) PRIMARY KEY,
			pet_id VARCHAR(36) NOT NULL,
			time VARCHAR(10) NOT NULL,
			food_type VARCHAR(50) DEFAULT '粮食',
			amount VARCHAR(50) DEFAULT '一份',
			FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS feeding_records (
			id VARCHAR(36) PRIMARY KEY,
			pet_id VARCHAR(36) NOT NULL,
			schedule_id VARCHAR(36) DEFAULT NULL,
			time VARCHAR(10) DEFAULT '',
			food_type VARCHAR(50) DEFAULT '粮食',
			amount VARCHAR(50) DEFAULT '一份',
			notes TEXT,
			created_at BIGINT NOT NULL,
			FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	}

	for _, ddl := range tables {
		if _, err := DB.Exec(ddl); err != nil {
			log.Fatalf("建表失败: %v", err)
		}
	}
	log.Println("数据库表初始化完成")
}
