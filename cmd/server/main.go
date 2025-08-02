// File: cmd/server/main.go
package main

import (
	"PanshiCMS/internal/config"   // 引入config包
	"PanshiCMS/internal/database" // 引入database包
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 1. 初始化配置
	config.InitConfig()

	// 2. 初始化数据库
	database.InitDB()

	// 3. 创建 Gin 引擎
	r := gin.Default()

	// 4. 创建路由处理器
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to PanshiCMS! Database connected.",
		})
	})

	// 5. 启动 HTTP 服务
	port := viper.GetString("server.port")
	if port == "" {
		port = ":8080" // 默认端口
	}
	r.Run(port)
}
