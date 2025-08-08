// File: internal/handler/ai_handler.go
package handler

import (
	"net/http"
	"strings"

	"PanshiCMS/internal/database"
	"PanshiCMS/internal/model"
	"github.com/gin-gonic/gin"
)

// GenerateMeta 通过调用大模型接口自动生成 SEO 元数据（标题和描述）。
// 该接口期望接收 JSON 参数：{"title":"标题","content":"文章正文"}。
// 如果站点设置中未配置大模型接入信息，则简单地根据传入内容截取生成默认元信息。
func GenerateMeta(c *gin.Context) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数无效"})
		return
	}
	// 读取站点设置，检查 AI 配置
	var settings model.SiteSettings
	database.DB.First(&settings)
	// 默认实现：如果未配置 AiApiKey，则直接截取内容作为元信息
	// 实际部署时，可根据 settings.AiAppId/AiApiKey/AiApiSecret 调用阿里云千问接口
	metaTitle := strings.TrimSpace(req.Title)
	metaDescription := ""
	if metaTitle == "" {
		// 如果没有标题，则取内容前 30 个字符作为标题
		if len(req.Content) > 30 {
			metaTitle = strings.TrimSpace(req.Content[:30])
		} else {
			metaTitle = strings.TrimSpace(req.Content)
		}
	}
	// 描述取内容前 60 个字符
	if len(req.Content) > 60 {
		metaDescription = strings.TrimSpace(req.Content[:60])
	} else {
		metaDescription = strings.TrimSpace(req.Content)
	}
	// TODO: 调用外部大模型 API 生成更智能的标题和描述。下面为示例伪代码：
	// if settings.AiApiKey != "" && settings.AiAppId != "" && settings.AiApiSecret != "" {
	//     response, err := service.CallQwenGenerateMeta(settings, req.Title, req.Content)
	//     if err == nil {
	//         metaTitle = response.Title
	//         metaDescription = response.Description
	//     }
	// }
	c.JSON(http.StatusOK, gin.H{
		"metaTitle":       metaTitle,
		"metaDescription": metaDescription,
	})
}
