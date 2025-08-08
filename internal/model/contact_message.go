// File: internal/model/contact_message.go
package model

import "gorm.io/gorm"

// ContactMessage 代表从前台“联系我们”表单提交的留言
// 该模型用于记录用户的姓名、邮箱、电话以及留言内容
type ContactMessage struct {
	gorm.Model
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Message string `gorm:"type:text" json:"message"`
}
