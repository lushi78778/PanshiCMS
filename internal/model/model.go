// File: internal/model/models.go
package model

import (
	"gorm.io/gorm"
	"time"
)

// NewsArticle 代表一条新闻或行业动态
type NewsArticle struct {
	gorm.Model

	Title         string    `gorm:"not null" json:"title"`
	CoverImageURL string    `json:"coverImageUrl"` // 改为小驼峰
	Summary       string    `json:"summary"`
	Content       string    `gorm:"type:text" json:"content"`
	PublishDate   time.Time `json:"publishDate"`                     // 改为小驼峰
	IsPublished   bool      `gorm:"default:true" json:"isPublished"` // 改为小驼峰

	Slug            string `gorm:"unique;not null" json:"slug"`
	MetaTitle       string `json:"metaTitle"`       // 改为小驼峰
	MetaDescription string `json:"metaDescription"` // 改为小驼峰
	MetaKeywords    string `json:"metaKeywords"`    // 改为小驼峰
}

// SiteSettings 代表网站的全局设置
type SiteSettings struct {
	gorm.Model

	CompanyName   string `json:"companyName"`
	ContactPhone  string `json:"contactPhone"`
	ContactEmail  string `json:"contactEmail"`
	Address       string `json:"address"`
	WechatQRURL   string `json:"wechatQrUrl"`
	ICPNumber     string `json:"icpNumber"`
	GonganBeian   string `json:"gonganBeian"`
	CopyrightInfo string `json:"copyrightInfo"`
}

// Service 代表一项产品或服务
type Service struct {
	gorm.Model

	Title         string `gorm:"not null" json:"title"`
	CoverImageURL string `json:"coverImageUrl"`
	Summary       string `gorm:"type:text" json:"summary"`
	Content       string `gorm:"type:text" json:"content"`
	IsRecommended bool   `gorm:"default:false" json:"isRecommended"`
	SortOrder     int    `gorm:"default:99" json:"sortOrder"`

	Slug            string `gorm:"unique;not null" json:"slug"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
}

// CaseStudy 代表一个成功案例
type CaseStudy struct {
	gorm.Model

	Title         string `gorm:"not null" json:"title"`
	CustomerName  string `json:"customerName"`
	Industry      string `json:"industry"`
	CoverImageURL string `json:"coverImageUrl"`
	Content       string `gorm:"type:text" json:"content"`
	IsPublished   bool   `gorm:"default:true" json:"isPublished"`

	Slug string `gorm:"unique;not null" json:"slug"`
}

// AdminUser 代表后台管理员用户
type AdminUser struct {
	gorm.Model
	Username     string `gorm:"unique;not null" json:"username"`
	PasswordHash string `gorm:"not null" json:"-"` // 在json输出时忽略密码哈希，增加安全性
}
