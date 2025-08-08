// File: internal/handler/public_handler.go
package handler

import (
	"fmt"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"

	"PanshiCMS/internal/database"
	"PanshiCMS/internal/model"
	"github.com/gin-gonic/gin"
)

// ShowHomePage 渲染首页
func ShowHomePage(c *gin.Context) {
	var settings model.SiteSettings
	database.DB.First(&settings)
	// 1. 查询首页显示的服务数量，根据站点设置的 HomepageServiceCount，如果为0则默认为6
	limit := settings.HomepageServiceCount
	if limit <= 0 {
		limit = 6
	}
	var services []model.Service
	database.DB.Order("is_recommended desc, sort_order asc").Limit(limit).Find(&services)

	var cases []model.CaseStudy
	database.DB.Order("created_at desc").Limit(6).Find(&cases)

	var news []model.NewsArticle
	database.DB.Where("is_published = ?", true).Order("publish_date desc").Limit(6).Find(&news)

	// 2. 查询发布的轮播图
	var banners []model.Banner
	database.DB.Where("is_published = ?", true).Order("sort_order asc, id asc").Find(&banners)
	// 如果没有任何轮播图，则回退到旧的设置图片或推荐服务封面
	if len(banners) == 0 {
		// 尝试从 site settings 中的三个banner字段读取
		if settings.BannerImage1URL != "" || settings.BannerImage2URL != "" || settings.BannerImage3URL != "" {
			if settings.BannerImage1URL != "" {
				banners = append(banners, model.Banner{ImageURL: settings.BannerImage1URL})
			}
			if settings.BannerImage2URL != "" {
				banners = append(banners, model.Banner{ImageURL: settings.BannerImage2URL})
			}
			if settings.BannerImage3URL != "" {
				banners = append(banners, model.Banner{ImageURL: settings.BannerImage3URL})
			}
		} else {
			// 最后回退到推荐服务封面
			for i, svc := range services {
				if i >= 3 {
					break
				}
				if svc.CoverImageURL != "" {
					banners = append(banners, model.Banner{ImageURL: svc.CoverImageURL})
				}
			}
		}
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"services": services,
		"cases":    cases,
		"news":     news,
		"settings": settings,
		"banners":  banners,
	})
}

// ShowPublicNewsList 渲染新闻列表页
func ShowPublicNewsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize := 10
	offset := (page - 1) * pageSize
	keyword := strings.TrimSpace(c.DefaultQuery("q", ""))
	var news []model.NewsArticle
	db := database.DB.Model(&model.NewsArticle{}).Where("is_published = ?", true)
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("title LIKE ? OR summary LIKE ?", like, like)
	}
	db.Order("publish_date desc").Limit(pageSize).Offset(offset).Find(&news)
	var settings model.SiteSettings
	database.DB.First(&settings)
	c.HTML(http.StatusOK, "news_list.html", gin.H{
		"news":     news,
		"page":     page,
		"keyword":  keyword,
		"settings": settings,
	})
}

// ShowPublicNewsDetail 渲染新闻详情页
func ShowPublicNewsDetail(c *gin.Context) {
	slug := c.Param("slug")
	var article model.NewsArticle
	if err := database.DB.Where("slug = ?", slug).First(&article).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	var settings model.SiteSettings
	database.DB.First(&settings)
	c.HTML(http.StatusOK, "news_detail.html", gin.H{
		"article":  article,
		"settings": settings,
	})
}

// ShowAboutPage 渲染关于我们页
func ShowAboutPage(c *gin.Context) {
	var settings model.SiteSettings
	database.DB.First(&settings)
	c.HTML(http.StatusOK, "about.html", gin.H{
		"settings": settings,
	})
}

// ShowContactPage 渲染联系我们页
func ShowContactPage(c *gin.Context) {
	var settings model.SiteSettings
	database.DB.First(&settings)
	c.HTML(http.StatusOK, "contact.html", gin.H{
		"settings": settings,
	})
}

// SubmitContactMessage 处理联系我们表单提交
func SubmitContactMessage(c *gin.Context) {
	var form struct {
		Name    string `form:"name" json:"name" binding:"required"`
		Email   string `form:"email" json:"email"`
		Phone   string `form:"phone" json:"phone"`
		Message string `form:"message" json:"message" binding:"required"`
	}
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "提交数据无效"})
		return
	}
	msg := model.ContactMessage{
		Name:    form.Name,
		Email:   form.Email,
		Phone:   form.Phone,
		Message: form.Message,
	}
	if err := database.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存留言失败"})
		return
	}
	// 尝试发送邮件通知给管理员
	var settings model.SiteSettings
	database.DB.First(&settings)
	// 当配置了 SMTP 信息和管理员邮箱时才发送
	if settings.SmtpHost != "" && settings.SmtpUser != "" && settings.SmtpPassword != "" && settings.ContactEmail != "" {
		host := settings.SmtpHost
		port := settings.SmtpPort
		// 默认端口
		if port == 0 {
			port = 25
		}
		from := settings.SmtpFrom
		if from == "" {
			from = settings.SmtpUser
		}
		to := []string{settings.ContactEmail}
		subject := "网站留言通知"
		body := fmt.Sprintf("您有一条新的留言：\n姓名：%s\n电话：%s\n邮箱：%s\n内容：\n%s", msg.Name, msg.Phone, msg.Email, msg.Message)
		message := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", settings.ContactEmail, subject, body))
		auth := smtp.PlainAuth("", settings.SmtpUser, settings.SmtpPassword, host)
		addr := fmt.Sprintf("%s:%d", host, port)
		// 忽略发送错误，避免影响用户提交体验
		_ = smtp.SendMail(addr, auth, from, to, message)
	}
	c.JSON(http.StatusOK, gin.H{"message": "感谢您的留言，我们会尽快联系您！"})
}

// ShowServicesPage 渲染产品/服务页面（当前为占位符）

// ShowServicesPage 渲染产品服务列表页面
func ShowServicesPage(c *gin.Context) {
	var settings model.SiteSettings
	database.DB.First(&settings)
	var services []model.Service
	database.DB.Order("sort_order asc, id asc").Find(&services)
	c.HTML(http.StatusOK, "services.html", gin.H{
		"settings": settings,
		"services": services,
	})
}

// ShowServiceDetail 渲染单个服务详情页
func ShowServiceDetail(c *gin.Context) {
	slug := c.Param("slug")
	var svc model.Service
	if err := database.DB.Where("slug = ?", slug).First(&svc).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	var settings model.SiteSettings
	database.DB.First(&settings)
	c.HTML(http.StatusOK, "service_detail.html", gin.H{
		"service":  svc,
		"settings": settings,
	})
}

// ShowCasesPage 渲染案例页面（当前为占位符）
// ShowCasesPage 渲染成功案例列表页面
func ShowCasesPage(c *gin.Context) {
	var settings model.SiteSettings
	database.DB.First(&settings)
	var cases []model.CaseStudy
	database.DB.Order("created_at desc").Find(&cases)
	c.HTML(http.StatusOK, "cases.html", gin.H{
		"settings": settings,
		"cases":    cases,
	})
}

// ShowCaseDetail 渲染案例详情页
func ShowCaseDetail(c *gin.Context) {
	slug := c.Param("slug")
	var cs model.CaseStudy
	if err := database.DB.Where("slug = ?", slug).First(&cs).Error; err != nil {
		c.HTML(http.StatusNotFound, "404.html", nil)
		return
	}
	var settings model.SiteSettings
	database.DB.First(&settings)
	c.HTML(http.StatusOK, "case_detail.html", gin.H{
		"case":     cs,
		"settings": settings,
	})
}

// ShowPrivacyPage 渲染隐私政策页面
func ShowPrivacyPage(c *gin.Context) {
	var settings model.SiteSettings
	database.DB.First(&settings)
	c.HTML(http.StatusOK, "privacy.html", gin.H{
		"settings": settings,
	})
}

// ShowLegalPage 渲染法律声明页面
func ShowLegalPage(c *gin.Context) {
	var settings model.SiteSettings
	database.DB.First(&settings)
	c.HTML(http.StatusOK, "legal.html", gin.H{
		"settings": settings,
	})
}
