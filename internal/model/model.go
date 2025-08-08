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

	// 合规性声明
	PrivacyPolicy  string `gorm:"type:text" json:"privacyPolicy"`  // 隐私政策正文，富文本HTML
	LegalStatement string `gorm:"type:text" json:"legalStatement"` // 法律声明正文，富文本HTML

	// 首页轮播图图片 URL，可配置三张，如有需要可拓展
	BannerImage1URL string `json:"bannerImage1Url"`
	BannerImage2URL string `json:"bannerImage2Url"`
	BannerImage3URL string `json:"bannerImage3Url"`

	// 首页服务推荐数量，控制首页显示多少个服务，默认6
	HomepageServiceCount int `json:"homepageServiceCount"`

	// 首页自定义介绍文本，可使用富文本格式
	HomepageIntro string `gorm:"type:text" json:"homepageIntro"`

	// 默认全局 SEO 标题和描述，可用于首页或缺省情况
	DefaultMetaTitle       string `json:"defaultMetaTitle"`
	DefaultMetaDescription string `json:"defaultMetaDescription"`

	// --- AI 与邮件配置 ---
	// AiEditor生成元信息时需要调用大模型 API。本系统支持配置大模型接入信息，
	// 由管理员在后台填写，对应阿里云通义千问服务的 AppID、API Key 和 Secret。
	AiAppId     string `json:"aiAppId"`     // 阿里云模型服务的 App ID
	AiApiKey    string `json:"aiApiKey"`    // API Key
	AiApiSecret string `json:"aiApiSecret"` // API Secret

	// 邮件通知配置：用于联系我们表单的邮件通知
	SmtpHost     string `json:"smtpHost"`     // 邮件服务器地址，如 smtpdm.aliyun.com
	SmtpPort     int    `json:"smtpPort"`     // 邮件服务器端口
	SmtpUser     string `json:"smtpUser"`     // SMTP 帐号
	SmtpPassword string `json:"smtpPassword"` // SMTP 密钥
	SmtpFrom     string `json:"smtpFrom"`     // 发件人邮箱地址（例如 info@lgic.top）
}

// Banner 代表首页轮播图数据
type Banner struct {
	gorm.Model
	ImageURL    string `json:"imageUrl"`
	Link        string `json:"link"`
	SortOrder   int    `gorm:"default:0" json:"sortOrder"`
	IsPublished bool   `gorm:"default:true" json:"isPublished"`
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
	// MetaKeywords 用于SEO关键词，逗号分隔多个关键词
	MetaKeywords string `json:"metaKeywords"`
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

	// SEO 字段
	MetaTitle       string `json:"metaTitle"`       // 页面标题
	MetaDescription string `json:"metaDescription"` // 页面描述
	MetaKeywords    string `json:"metaKeywords"`    // 关键词
}

// AdminUser 代表后台管理员用户
type AdminUser struct {
	gorm.Model
	Username     string `gorm:"unique;not null" json:"username"`
	PasswordHash string `gorm:"not null" json:"-"` // 在json输出时忽略密码哈希，增加安全性
}
