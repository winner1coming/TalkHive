package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config 用于存储应用的整体配置
type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dsn          string
		MaxIdleConns int
		MaxOpenConns int
	}
	Redis struct {
		Addr     string
		DB       int
		Password string
	}
}

// AppConfig 全局配置变量
var AppConfig *Config

// InitConfig 初始化配置文件
func InitConfig() {
	// 设置配置文件名称、类型和路径
	viper.SetConfigName("config")   // 配置文件名称为 config
	viper.SetConfigType("yaml")     // 配置文件类型为 yaml
	viper.AddConfigPath("./config") // 配置文件路径

	// 读取并解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("配置文件读取失败: %v", err)
	}

	// 创建配置实例
	AppConfig = &Config{}

	// 将配置文件解析到 AppConfig
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("配置文件解析失败: %v", err)
	}

	log.Println("配置文件解析成功")
}
