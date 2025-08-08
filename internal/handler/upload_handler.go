// File: internal/handler/upload_handler.go
package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadFile 处理文件上传请求
// 前端应以 multipart/form-data 格式提交文件，字段名为 "file"
// 上传后的文件将保存到 web/static/uploads 目录下，返回值包含可访问的相对路径
func UploadFile(c *gin.Context) {
	// 从表单获取文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到上传的文件"})
		return
	}

	// 构建保存路径并确保目录存在
	uploadDir := filepath.Join("web", "static", "uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
		return
	}

	// 生成唯一文件名，避免重复
	ext := filepath.Ext(file.Filename)
	newName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dstPath := filepath.Join(uploadDir, newName)

	// 保存文件
	if err := c.SaveUploadedFile(file, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}

	// 返回上传后的URL路径
	fileURL := "/uploads/" + newName
	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
		"url":     fileURL,
	})
}
