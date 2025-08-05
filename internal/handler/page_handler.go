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
