// File: cmd/server/main.go
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的 Gin 引擎
	r := gin.Default()

	// 创建一个简单的路由处理器
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to PanshiCMS!",
		})
	})

	// 启动 HTTP 服务，默认在 0.0.0.0:8080 启动
	r.Run()
}
