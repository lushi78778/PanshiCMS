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

	// 提供静态文件服务，用于上传的图片和前端静态资源
	// 上传的文件统一暴露在 /uploads 路径下：
	// 例如，上传的图片存放在 web/static/uploads/xxx.png，则可以通过 /uploads/xxx.png 访问到
	r.Static("/uploads", filepath.Join("web", "static", "uploads"))
	// 前端静态资源（如 AiEditor）的访问路径：
	// 任何放在 web/static 中的文件可以通过 /static 访问到，例如
	// web/static/aieditor/index.js -> /static/aieditor/index.js
	r.Static("/static", filepath.Join("web", "static"))

	// --- 前台页面路由 ---
	// sitemap.xml 用于 SEO
	r.GET("/sitemap.xml", handler.SitemapHandler)
	r.GET("/", handler.ShowHomePage)
	r.GET("/about", handler.ShowAboutPage)
	r.GET("/news", handler.ShowPublicNewsList)
	r.GET("/news/:slug", handler.ShowPublicNewsDetail)
	r.GET("/contact", handler.ShowContactPage)
	r.POST("/contact", handler.SubmitContactMessage)
	r.GET("/services", handler.ShowServicesPage)
	r.GET("/services/:slug", handler.ShowServiceDetail)
	r.GET("/cases", handler.ShowCasesPage)
	r.GET("/cases/:slug", handler.ShowCaseDetail)

	// 合规性页面
	r.GET("/privacy", handler.ShowPrivacyPage)
	r.GET("/legal", handler.ShowLegalPage)

	// --- 后台页面路由 (供浏览器访问) ---
	adminPageRoutes := r.Group("/admin")
	{
		adminPageRoutes.GET("/login", handler.ShowLoginPage)
		adminPageRoutes.GET("/dashboard", handler.ShowAdminDashboard)
		adminPageRoutes.GET("/news", handler.ShowAdminNewsList)
		adminPageRoutes.GET("/news/new", handler.ShowAdminNewsEditPage)      // <-- 新增：新建文章页面
		adminPageRoutes.GET("/news/edit/:id", handler.ShowAdminNewsEditPage) // <-- 新增：编辑文章页面

		// 服务管理页面
		adminPageRoutes.GET("/services", handler.ShowAdminServicesList)
		adminPageRoutes.GET("/services/new", handler.ShowAdminServiceEditPage)
		adminPageRoutes.GET("/services/edit/:id", handler.ShowAdminServiceEditPage)

		// 案例管理页面
		adminPageRoutes.GET("/cases", handler.ShowAdminCasesList)
		adminPageRoutes.GET("/cases/new", handler.ShowAdminCaseEditPage)
		adminPageRoutes.GET("/cases/edit/:id", handler.ShowAdminCaseEditPage)

		// 轮播图管理页面
		adminPageRoutes.GET("/banners", handler.ShowAdminBannersList)
		adminPageRoutes.GET("/banners/new", handler.ShowAdminBannerEditPage)
		adminPageRoutes.GET("/banners/edit/:id", handler.ShowAdminBannerEditPage)

		// 站点设置页面
		adminPageRoutes.GET("/settings", handler.ShowAdminSettingsPage)
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

		// 文件上传API
		// 用于处理后台富文本编辑器和封面上传等需求
		apiV1.POST("/upload", handler.UploadFile)

		// 站点设置API
		apiV1.GET("/site-settings", handler.GetSiteSettings)
		apiV1.PUT("/site-settings", handler.UpdateSiteSettings)

		// AI 元数据生成接口
		apiV1.POST("/ai/generate-meta", handler.GenerateMeta)

		// 轮播图API
		bannerRoutes := apiV1.Group("/banners")
		{
			bannerRoutes.POST("/", handler.CreateBanner)
			bannerRoutes.GET("/", handler.GetBannerList)
			bannerRoutes.GET("/:id", handler.GetBannerByID)
			bannerRoutes.PUT("/:id", handler.UpdateBanner)
			bannerRoutes.DELETE("/:id", handler.DeleteBanner)
		}
	}

	// 兼容旧约定的初始化入口 `/do?action=init`
	// 此路由用于在系统未初始化时创建第一个管理员账户
	// 访问 /do?action=init 将触发 InitAdminHandler
	r.GET("/do", handler.DoHandler)
}

// loadTemplates 是一个健壮的辅助函数，用于加载指定目录下的所有模板
func loadTemplates(templatesDir string) *template.Template {
	// 1. 创建一个主模板对象，这是所有模板的集合
	//    同时注册自定义函数，如 safe 用于输出 HTML 内容
	tmpl := template.New("").Funcs(template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	})

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
