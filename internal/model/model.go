// File: internal/model/models.go
package model

import (
	"gorm.io/gorm"
	"time"
)

// NewsArticle 代表一条新闻或行业动态
type NewsArticle struct {
	gorm.Model // 内嵌gorm.Model，自带ID, CreatedAt, UpdatedAt, DeletedAt

	Title         string    `gorm:"not null" json:"title"`            // 标题，不能为空
	CoverImageURL string    `json:"cover_image_url"`                  // 封面图URL
	Summary       string    `json:"summary"`                          // 摘要
	Content       string    `gorm:"type:text" json:"content"`         // 正文，富文本HTML
	PublishDate   time.Time `json:"publish_date"`                     // 发布日期
	IsPublished   bool      `gorm:"default:true" json:"is_published"` // 是否发布状态

	// SEO 字段
	Slug            string `gorm:"unique;not null" json:"slug"` // URL别名，用于生成SEF URL，唯一且不能为空
	MetaTitle       string `json:"meta_title"`                  // SEO标题（留空则默认使用主标题）
	MetaDescription string `json:"meta_description"`            // SEO描述
	MetaKeywords    string `json:"meta_keywords"`               // SEO关键词
}

// SiteSettings 代表网站的全局设置，这是一个单例模型，表中通常只有一条记录
type SiteSettings struct {
	gorm.Model

	CompanyName   string `json:"company_name"`
	ContactPhone  string `json:"contact_phone"`
	ContactEmail  string `json:"contact_email"`
	Address       string `json:"address"`
	WechatQRURL   string `json:"wechat_qr_url"`
	ICPNumber     string `json:"icp_number"`
	GonganBeian   string `json:"gongan_beian"`
	CopyrightInfo string `json:"copyright_info"`
	// ... 未来可以继续添加其他全局设置 ...
}

// Service 代表一项产品或服务
type Service struct {
	gorm.Model

	Title         string `gorm:"not null" json:"title"`
	CoverImageURL string `json:"cover_image_url"`
	Summary       string `gorm:"type:text" json:"summary"`
	Content       string `gorm:"type:text" json:"content"`
	IsRecommended bool   `gorm:"default:false" json:"is_recommended"` // 是否推荐到首页轮播图
	SortOrder     int    `gorm:"default:99" json:"sort_order"`        // 排序值

	// SEO 字段
	Slug            string `gorm:"unique;not null" json:"slug"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
}

// CaseStudy 代表一个成功案例
type CaseStudy struct {
	gorm.Model

	Title         string `gorm:"not null" json:"title"`
	CustomerName  string `json:"customer_name"`
	Industry      string `json:"industry"`
	CoverImageURL string `json:"cover_image_url"`
	Content       string `gorm:"type:text" json:"content"`
	IsPublished   bool   `gorm:"default:true" json:"is_published"`

	// SEO 字段
	Slug string `gorm:"unique;not null" json:"slug"`
}

// AdminUser 代表后台管理员用户
type AdminUser struct {
	gorm.Model
	Username     string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	// 可选：未来可增加角色、权限等字段
}
