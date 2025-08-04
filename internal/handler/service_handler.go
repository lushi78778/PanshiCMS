// File: internal/handler/service_handler.go
package handler

import (
	"PanshiCMS/internal/database"
	"PanshiCMS/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateService 创建产品/服务
func CreateService(c *gin.Context) {
	var service model.Service
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	if err := database.DB.Create(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建服务失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "服务创建成功", "data": service})
}

// GetServiceList 获取产品/服务列表
func GetServiceList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	var serviceList []model.Service
	var total int64

	database.DB.Model(&model.Service{}).Count(&total)
	database.DB.Order("sort_order asc, created_at desc").Limit(pageSize).Offset(offset).Find(&serviceList)

	c.JSON(http.StatusOK, gin.H{
		"message": "获取服务列表成功",
		"data": gin.H{
			"list":     serviceList,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetServiceByID 根据ID获取单个产品/服务
func GetServiceByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var service model.Service
	if err := database.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务未找到"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取服务成功", "data": service})
}

// UpdateService 更新产品/服务
func UpdateService(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var service model.Service
	if err := database.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "服务未找到"})
		return
	}
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	if err := database.DB.Save(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新服务失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "服务更新成功", "data": service})
}

// DeleteService 删除产品/服务
func DeleteService(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := database.DB.Delete(&model.Service{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除服务失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "服务删除成功"})
}
