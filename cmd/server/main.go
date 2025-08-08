// File: cmd/server/main.go
package main

import (
	"PanshiCMS/internal/config"
	"PanshiCMS/internal/database"
	"PanshiCMS/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.InitConfig()
	database.InitDB()

	r := gin.Default()

	// 设置路由
	router.SetupRouter(r)

	// 启动 HTTP 服务
	port := viper.GetString("server.port")
	if port == "" {
		port = ":8080"
	}
	r.Run(port)
}
