package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB     DBConfig `yaml:"db"`
	Server Server   `yaml:"server"`
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
	path := filepath.Join("config", "config.yaml")
	if v := os.Getenv("CONFIG_PATH"); v != "" {
		path = v
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return defaultConfig()
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return defaultConfig()
	}
	return &cfg
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
