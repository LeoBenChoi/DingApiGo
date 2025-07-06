package config

import (
	"log"
	"os"

	// 使用 yaml.v3
	"gopkg.in/yaml.v3"
)

// DingDingConfig 钉钉配置结构体
type DingDingConfig struct {
	AppKey    string `yaml:"app_key"`    // 钉钉应用的 AppKey
	AppSecret string `yaml:"app_secret"` // 钉钉应用的 AppSecret
	AgentID   int64  `yaml:"agent_id"`   // 钉钉应用的 AgentID，通常为 int64 类型
}

// Config 全局配置结构体
type Config struct {
	DingDing DingDingConfig `yaml:"dingding"` // 钉钉相关配置
}

var cfg Config

// LoadConfig 加载并解析 config.yaml 配置文件
func LoadConfig() {
	// 读取配置文件
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 解析 YAML 数据到 Config 结构体
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}
}

// GetConfig 获取全局配置
func GetConfig() Config {
	return cfg
}
