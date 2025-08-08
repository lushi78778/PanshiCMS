// File: internal/handler/page_handler.go
package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ShowLoginPage 渲染登录页面
func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/login.html", nil)
}

// ShowAdminDashboard 渲染后台主页面
func ShowAdminDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/dashboard.html", nil)
}

// ShowAdminNewsList 渲染后台新闻列表页面
func ShowAdminNewsList(c *gin.Context) {
	// 将之前的重定向，修改为渲染我们新的模板文件
	c.HTML(http.StatusOK, "admin/news_list.html", nil)
}

// ShowAdminNewsEditPage 渲染新闻的新建/编辑页面
func ShowAdminNewsEditPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/news_edit.html", nil)
}

// ShowAdminServicesList 渲染后台服务列表页面
func ShowAdminServicesList(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/services_list.html", nil)
}

// ShowAdminServiceEditPage 渲染后台服务新建/编辑页面
func ShowAdminServiceEditPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/services_edit.html", nil)
}

// ShowAdminCasesList 渲染后台案例列表页面
func ShowAdminCasesList(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/cases_list.html", nil)
}

// ShowAdminCaseEditPage 渲染后台案例新建/编辑页面
func ShowAdminCaseEditPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/cases_edit.html", nil)
}

// ShowAdminSettingsPage 渲染站点设置编辑页面
func ShowAdminSettingsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/settings.html", nil)
}

// ShowAdminBannersList 渲染轮播图列表页面
func ShowAdminBannersList(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/banners_list.html", nil)
}

// ShowAdminBannerEditPage 渲染轮播图编辑/新建页面
func ShowAdminBannerEditPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/banner_edit.html", nil)
}
