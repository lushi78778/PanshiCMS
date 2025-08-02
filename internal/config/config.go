// File: internal/config/config.go
package config

import (
	"github.com/spf13/viper"
	"log"
)

// InitConfig 初始化Viper配置
func InitConfig() {
	viper.SetConfigName("config") // 配置文件名 (不带扩展名)
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(".")      // 在当前目录查找配置文件

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalln("找不到配置文件 'config.yaml'")
		} else {
			log.Fatalf("读取配置文件失败: %v", err)
		}
	}
	log.Println("配置文件加载成功!")
}
