// File: internal/handler/banner_handler.go
package handler

import (
	"net/http"
	"strconv"

	"PanshiCMS/internal/database"
	"PanshiCMS/internal/model"
	"github.com/gin-gonic/gin"
)

// CreateBanner 创建新的轮播图条目
func CreateBanner(c *gin.Context) {
	var input model.Banner
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据无效: " + err.Error()})
		return
	}
	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建轮播图失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "data": input})
}

// GetBannerList 获取轮播图列表
func GetBannerList(c *gin.Context) {
	var banners []model.Banner
	database.DB.Order("sort_order asc, id asc").Find(&banners)
	c.JSON(http.StatusOK, gin.H{"data": banners})
}

// GetBannerByID 获取单个轮播图
func GetBannerByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID 无效"})
		return
	}
	var banner model.Banner
	if err := database.DB.First(&banner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": banner})
}

// UpdateBanner 更新轮播图
func UpdateBanner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID 无效"})
		return
	}
	var banner model.Banner
	if err := database.DB.First(&banner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}
	var input model.Banner
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}
	banner.ImageURL = input.ImageURL
	banner.Link = input.Link
	banner.SortOrder = input.SortOrder
	banner.IsPublished = input.IsPublished
	if err := database.DB.Save(&banner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": banner})
}

// DeleteBanner 删除轮播图
func DeleteBanner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID 无效"})
		return
	}
	if err := database.DB.Delete(&model.Banner{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
