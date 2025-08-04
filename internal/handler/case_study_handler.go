// File: internal/handler/case_study_handler.go
package handler

import (
	"PanshiCMS/internal/database"
	"PanshiCMS/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateCaseStudy 创建成功案例
func CreateCaseStudy(c *gin.Context) {
	var caseStudy model.CaseStudy
	if err := c.ShouldBindJSON(&caseStudy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	if err := database.DB.Create(&caseStudy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建案例失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "案例创建成功", "data": caseStudy})
}

// GetCaseStudyList 获取成功案例列表
func GetCaseStudyList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	var caseStudyList []model.CaseStudy
	var total int64

	database.DB.Model(&model.CaseStudy{}).Count(&total)
	database.DB.Order("created_at desc").Limit(pageSize).Offset(offset).Find(&caseStudyList)

	c.JSON(http.StatusOK, gin.H{
		"message": "获取案例列表成功",
		"data": gin.H{
			"list":     caseStudyList,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetCaseStudyByID 根据ID获取单个成功案例
func GetCaseStudyByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var caseStudy model.CaseStudy
	if err := database.DB.First(&caseStudy, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "案例未找到"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取案例成功", "data": caseStudy})
}

// UpdateCaseStudy 更新成功案例
func UpdateCaseStudy(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var caseStudy model.CaseStudy
	if err := database.DB.First(&caseStudy, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "案例未找到"})
		return
	}
	if err := c.ShouldBindJSON(&caseStudy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	if err := database.DB.Save(&caseStudy).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新案例失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "案例更新成功", "data": caseStudy})
}

// DeleteCaseStudy 删除成功案例
func DeleteCaseStudy(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := database.DB.Delete(&model.CaseStudy{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除案例失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "案例删除成功"})
}
