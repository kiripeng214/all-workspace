package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB     DBConfig    `yaml:"db"`
	Server Server      `yaml:"server"`
	LLM    LLMConfig   `yaml:"llm"`
}

type LLMConfig struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
	APIURL   string `yaml:"api_url"`
	Model    string `yaml:"model"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Server struct {
	Port string `yaml:"port"`
}

func Load() *Config {
	// 1. 加载主配置
	path := filepath.Join("config", "config.yaml")
	if v := os.Getenv("CONFIG_PATH"); v != "" {
		path = v
	}

	cfg := defaultConfig()

	data, err := os.ReadFile(path)
	if err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			log.Printf("⚠️ 配置解析失败，使用默认值: %v", err)
		}
	}

	// 2. 加载本地覆盖配置（不上传）
	localPath := filepath.Join("config", "config-local.yaml")
	if v := os.Getenv("CONFIG_LOCAL_PATH"); v != "" {
		localPath = v
	}
	if localData, err := os.ReadFile(localPath); err == nil {
		if err := yaml.Unmarshal(localData, cfg); err != nil {
			log.Printf("⚠️ 本地配置解析失败: %v", err)
		}
	}

	return cfg
}

func defaultConfig() *Config {
	return &Config{
		DB: DBConfig{
			Host: "localhost",
			Port: "3306",
			User: "root",
			Name: "pet_applet",
		},
		Server: Server{
			Port: "3000",
		},
	}
}
