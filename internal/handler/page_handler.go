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

// ShowAdminNewsList 渲染后台新闻列表页面 (为下一步做准备)
func ShowAdminNewsList(c *gin.Context) {
	// 暂时先重定向到dashboard，之后会实现具体页面
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

// ... 未来会添加更多页面渲染函数 ...
