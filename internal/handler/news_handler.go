// File: internal/handler/news_handler.go
package handler

import (
	"net/http"
	"strconv"
	"strings"

	"PanshiCMS/internal/database"
	"PanshiCMS/internal/model"
	"github.com/gin-gonic/gin"
)

// CreateNews 创建一篇新文章
func CreateNews(c *gin.Context) {
	var news model.NewsArticle
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	if err := database.DB.Create(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章创建成功", "data": news})
}

// GetNewsList 获取文章列表（含分页）
func GetNewsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	// 支持简单的标题/摘要搜索，通过查询参数 q 传入关键字
	keyword := strings.TrimSpace(c.DefaultQuery("q", ""))

	var newsList []model.NewsArticle
	var total int64

	// 构建基础查询
	db := database.DB.Model(&model.NewsArticle{})
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("title LIKE ? OR summary LIKE ?", like, like)
	}

	// 获取总数
	db.Count(&total)
	// 查询分页数据
	db.Order("publish_date desc").Limit(pageSize).Offset(offset).Find(&newsList)

	c.JSON(http.StatusOK, gin.H{
		"message": "获取文章列表成功",
		"data": gin.H{
			"list":     newsList,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetNewsByID 根据ID获取单篇文章
func GetNewsByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var news model.NewsArticle
	if err := database.DB.First(&news, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "获取文章成功", "data": news})
}

// UpdateNews 更新一篇文章
func UpdateNews(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var news model.NewsArticle
	if err := database.DB.First(&news, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}

	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}

	if err := database.DB.Save(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章更新成功", "data": news})
}

// DeleteNews 删除一篇文章
func DeleteNews(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := database.DB.Delete(&model.NewsArticle{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}
