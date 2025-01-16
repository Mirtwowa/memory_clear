package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// Config 结构体用于存储配置文件中的所有配置项
type Config struct {
	MemoryThreshold   int                `yaml:"memorythreshold"`   // 内存占用阈值（单位：百分比）
	ProcessIgnoreList []string           `yaml:"processignorelist"` // 忽略的进程列表
	Notification      NotificationConfig `yaml:"notification"`      // 通知设置
}

// NotificationConfig 存储通知相关的配置
type NotificationConfig struct {
	Enabled bool          `yaml:"enabled"` // 是否启用通知
	Method  string        `yaml:"method"`  // 通知方式（email 或 desktop）
	Email   EmailConfig   `yaml:"email"`
	Desktop DesktopConfig `yaml:"desktop"`
}

// EmailConfig 邮件通知配置
type EmailConfig struct {
	Recipient  string `yaml:"recipient"`   // 邮件接收者
	SMTPServer string `yaml:"smtp_server"` // SMTP 服务器地址
}

// DesktopConfig 桌面通知配置
type DesktopConfig struct {
	Enabled bool `yaml:"enabled"` // 是否启用桌面通知
	Timeout int  `yaml:"timeout"` // 通知显示时间（秒）
}

// configInstance 用于存储配置数据
var configInstance *Config

// GetConfig 返回配置实例
func GetConfig() *Config {
	// 如果 configInstance 为空，初始化配置
	if configInstance == nil {
		configInstance = &Config{}
		err := loadConfig(configInstance)
		if err != nil {
			log.Fatalf("加载配置失败: %v", err)
		}
	}
	return configInstance
}

// loadConfig 读取配置文件并将其解析到 Config 结构体
func loadConfig(c *Config) error {
	viper.SetConfigName("config")   // 配置文件名 (不带扩展名)
	viper.AddConfigPath("./config") // 配置文件路径
	viper.AutomaticEnv()            // 自动从环境变量加载配置

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("无法读取配置文件: %v", err)
	}

	// 解析配置
	if err := viper.Unmarshal(c); err != nil {
		return fmt.Errorf("无法解析配置文件: %v", err)
	}

	// 打印读取的配置，检查是否正确加载
	fmt.Printf("配置加载成功：%+v\n", *c)
	return nil
}
