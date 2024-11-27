package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

// Config 结构体用于存储配置信息
type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

// DatabaseConfig 存储数据库连接配置信息
type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

// LoadConfig 从 YAML 文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, err
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// 解析 YAML 内容
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
