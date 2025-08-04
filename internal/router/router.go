// File: internal/router/router.go
package router

import (
	"PanshiCMS/internal/handler"
	"PanshiCMS/internal/middleware" // 引入中间件包
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// 公共路由，无需认证
	r.POST("/init", handler.InitAdminHandler)    // 初始化管理员
	r.POST("/admin/login", handler.LoginHandler) // 管理员登录

	// API V1 路由组，需要认证
	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.AuthMiddleware()) // <--- 对整个v1路由组应用认证中间件
	{
		newsRoutes := apiV1.Group("/news")
		{
			newsRoutes.POST("/", handler.CreateNews)
			newsRoutes.GET("/", handler.GetNewsList)
			newsRoutes.GET("/:id", handler.GetNewsByID)
			newsRoutes.PUT("/:id", handler.UpdateNews)
			newsRoutes.DELETE("/:id", handler.DeleteNews)
		}
	}
}
