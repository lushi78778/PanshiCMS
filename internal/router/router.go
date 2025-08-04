// File: internal/router/router.go
package router

import (
	"PanshiCMS/internal/handler"
	"PanshiCMS/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// 公共API路由
	r.POST("/init", handler.InitAdminHandler)
	r.POST("/admin/login", handler.LoginHandler)

	// 受保护的API V1 路由组
	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.AuthMiddleware())
	{
		// 新闻文章路由组
		newsRoutes := apiV1.Group("/news")
		{
			newsRoutes.POST("/", handler.CreateNews)
			newsRoutes.GET("/", handler.GetNewsList)
			newsRoutes.GET("/:id", handler.GetNewsByID)
			newsRoutes.PUT("/:id", handler.UpdateNews)
			newsRoutes.DELETE("/:id", handler.DeleteNews)
		}

		// 产品/服务路由组
		serviceRoutes := apiV1.Group("/services")
		{
			serviceRoutes.POST("/", handler.CreateService)
			serviceRoutes.GET("/", handler.GetServiceList)
			serviceRoutes.GET("/:id", handler.GetServiceByID)
			serviceRoutes.PUT("/:id", handler.UpdateService)
			serviceRoutes.DELETE("/:id", handler.DeleteService)
		}

		// 成功案例路由组
		caseStudyRoutes := apiV1.Group("/cases")
		{
			caseStudyRoutes.POST("/", handler.CreateCaseStudy)
			caseStudyRoutes.GET("/", handler.GetCaseStudyList)
			caseStudyRoutes.GET("/:id", handler.GetCaseStudyByID)
			caseStudyRoutes.PUT("/:id", handler.UpdateCaseStudy)
			caseStudyRoutes.DELETE("/:id", handler.DeleteCaseStudy)
		}
	}
}
