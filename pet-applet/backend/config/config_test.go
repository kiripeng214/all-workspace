package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()
	if cfg == nil {
		t.Fatal("defaultConfig() 返回 nil")
	}
	if cfg.Server.Port != "3000" {
		t.Errorf("默认端口应为 3000，实际 %s", cfg.Server.Port)
	}
	if cfg.DB.Host != "localhost" {
		t.Errorf("默认数据库主机应为 localhost，实际 %s", cfg.DB.Host)
	}
	if cfg.DB.Name != "pet_applet" {
		t.Errorf("默认数据库名应为 pet_applet，实际 %s", cfg.DB.Name)
	}
}

func TestLoad_WithDefault(t *testing.T) {
	// 当配置文件不存在时，应返回默认配置
	orig := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", orig)
	os.Unsetenv("CONFIG_PATH")

	// 切换到不包含配置文件的临时目录
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	cfg := Load()
	if cfg.Server.Port != "3000" {
		t.Errorf("配置文件不存在时应返回默认配置，端口 %s", cfg.Server.Port)
	}
}

func TestLoad_WithCustomPath(t *testing.T) {
	// 创建临时配置文件
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configContent := []byte(`
server:
  port: "8080"
db:
  host: "192.168.1.1"
  port: "3306"
  user: "test_user"
  password: "test_pass"
  name: "test_db"
`)
	configPath := filepath.Join(tmpDir, "custom.yaml")
	if err := os.WriteFile(configPath, configContent, 0644); err != nil {
		t.Fatal(err)
	}

	orig := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", orig)
	os.Setenv("CONFIG_PATH", configPath)

	cfg := Load()
	if cfg.Server.Port != "8080" {
		t.Errorf("应加载自定义配置的端口 8080，实际 %s", cfg.Server.Port)
	}
	if cfg.DB.Host != "192.168.1.1" {
		t.Errorf("应加载自定义配置的数据库主机，实际 %s", cfg.DB.Host)
	}
	if cfg.DB.User != "test_user" {
		t.Errorf("应加载自定义配置的用户名，实际 %s", cfg.DB.User)
	}
	if cfg.DB.Name != "test_db" {
		t.Errorf("应加载自定义配置的数据库名，实际 %s", cfg.DB.Name)
	}
}

func TestLoad_InvalidYaml(t *testing.T) {
	// 配置文件损坏时应回退到默认配置
	tmpDir, err := os.MkdirTemp("", "config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "bad.yaml")
	if err := os.WriteFile(configPath, []byte(": invalid yaml :: "), 0644); err != nil {
		t.Fatal(err)
	}

	orig := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", orig)
	os.Setenv("CONFIG_PATH", configPath)

	cfg := Load()
	if cfg.Server.Port != "3000" {
		t.Errorf("YAML 损坏时应返回默认配置，端口 %s", cfg.Server.Port)
	}
}
