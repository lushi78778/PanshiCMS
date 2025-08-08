// File: internal/handler/settings_handler.go
package handler

import (
	"net/http"

	"PanshiCMS/internal/database"
	"PanshiCMS/internal/model"
	"github.com/gin-gonic/gin"
)

// GetSiteSettings 返回唯一的站点设置记录
func GetSiteSettings(c *gin.Context) {
	var settings model.SiteSettings
	if err := database.DB.First(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取站点设置失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": settings})
}

// UpdateSiteSettings 更新站点设置，包括隐私政策和法律声明
func UpdateSiteSettings(c *gin.Context) {
	var input model.SiteSettings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	var settings model.SiteSettings
	if err := database.DB.First(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取站点设置失败: " + err.Error()})
		return
	}
	// 更新需要修改的字段
	settings.CompanyName = input.CompanyName
	settings.ContactPhone = input.ContactPhone
	settings.ContactEmail = input.ContactEmail
	settings.Address = input.Address
	settings.WechatQRURL = input.WechatQRURL
	settings.ICPNumber = input.ICPNumber
	settings.GonganBeian = input.GonganBeian
	settings.CopyrightInfo = input.CopyrightInfo
	settings.PrivacyPolicy = input.PrivacyPolicy
	settings.LegalStatement = input.LegalStatement
	settings.BannerImage1URL = input.BannerImage1URL
	settings.BannerImage2URL = input.BannerImage2URL
	settings.BannerImage3URL = input.BannerImage3URL
	settings.HomepageServiceCount = input.HomepageServiceCount
	settings.HomepageIntro = input.HomepageIntro
	settings.DefaultMetaTitle = input.DefaultMetaTitle
	settings.DefaultMetaDescription = input.DefaultMetaDescription

	// 更新 AI 配置
	settings.AiAppId = input.AiAppId
	settings.AiApiKey = input.AiApiKey
	settings.AiApiSecret = input.AiApiSecret
	// 更新邮件配置
	settings.SmtpHost = input.SmtpHost
	settings.SmtpPort = input.SmtpPort
	settings.SmtpUser = input.SmtpUser
	settings.SmtpPassword = input.SmtpPassword
	settings.SmtpFrom = input.SmtpFrom

	if err := database.DB.Save(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新站点设置失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "站点设置更新成功", "data": settings})
}
