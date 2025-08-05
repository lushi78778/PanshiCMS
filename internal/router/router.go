// File: internal/router/router.go
package router

import (
	"PanshiCMS/internal/handler"
	"PanshiCMS/internal/middleware"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func SetupRouter(r *gin.Engine) {
	// --- 最终修正的模板加载方案 ---
	// 使用我们自定义的、支持嵌套目录的模板渲染器
	r.SetHTMLTemplate(loadTemplates("web/templates"))

	// --- 后台页面路由 (供浏览器访问) ---
	adminPageRoutes := r.Group("/admin")
	{
		adminPageRoutes.GET("/login", handler.ShowLoginPage)
		adminPageRoutes.GET("/dashboard", handler.ShowAdminDashboard)
		adminPageRoutes.GET("/news", handler.ShowAdminNewsList)
		adminPageRoutes.GET("/news/new", handler.ShowAdminNewsEditPage)      // <-- 新增：新建文章页面
		adminPageRoutes.GET("/news/edit/:id", handler.ShowAdminNewsEditPage) // <-- 新增：编辑文章页面
	}

	// --- 后台API路由 (供前端JS调用) ---
	r.POST("/init", handler.InitAdminHandler)
	r.POST("/admin/login", handler.LoginHandler)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.AuthMiddleware())
	{
		// 新闻API
		newsRoutes := apiV1.Group("/news")
		{
			newsRoutes.POST("/", handler.CreateNews)
			newsRoutes.GET("/", handler.GetNewsList)
			newsRoutes.GET("/:id", handler.GetNewsByID)
			newsRoutes.PUT("/:id", handler.UpdateNews)
			newsRoutes.DELETE("/:id", handler.DeleteNews)
		}
		// 服务API
		serviceRoutes := apiV1.Group("/services")
		{
			serviceRoutes.POST("/", handler.CreateService)
			serviceRoutes.GET("/", handler.GetServiceList)
			serviceRoutes.GET("/:id", handler.GetServiceByID)
			serviceRoutes.PUT("/:id", handler.UpdateService)
			serviceRoutes.DELETE("/:id", handler.DeleteService)
		}
		// 案例API
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

// loadTemplates 是一个健壮的辅助函数，用于加载指定目录下的所有模板
func loadTemplates(templatesDir string) *template.Template {
	// 1. 创建一个主模板对象，这是所有模板的集合
	tmpl := template.New("")

	// 2. 遍历模板目录
	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 3. 只处理.html文件
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			// 4. 获取相对路径作为模板的唯一名称
			name := strings.TrimPrefix(path, templatesDir+"/")

			// 5. 读取文件内容
			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("读取模板文件失败: %s, 错误: %v", path, err)
				return err
			}

			// 6. 创建一个新的、在主模板集合中有独立名字的子模板
			//    然后将文件内容(字符串)解析到这个子模板中
			_, err = tmpl.New(name).Parse(string(content))
			if err != nil {
				log.Printf("解析模板失败: %s, 错误: %v", path, err)
				return err
			}
			log.Printf("成功加载模板: %s", name)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return tmpl
}
