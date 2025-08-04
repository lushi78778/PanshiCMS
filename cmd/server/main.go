// File: cmd/server/main.go
package main

import (
	"PanshiCMS/internal/config"
	"PanshiCMS/internal/database"
	"PanshiCMS/internal/router" // 引入router包
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.InitConfig()
	database.InitDB()

	r := gin.Default()

	// 设置路由
	router.SetupRouter(r) // 调用我们的路由设置函数

	// 启动 HTTP 服务
	port := viper.GetString("server.port")
	if port == "" {
		port = ":8080"
	}
	r.Run(port)
}
