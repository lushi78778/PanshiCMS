// File: internal/database/database.go
package database

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"PanshiCMS/internal/model" // 请将 "your-repo" 替换为您的代码仓库地址
)

var DB *gorm.DB

// InitDB 根据配置初始化数据库连接
func InitDB() {
	var err error
	dbType := viper.GetString("database.type")

	// 配置GORM日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	gormConfig := &gorm.Config{
		Logger: newLogger,
	}

	switch dbType {
	case "sqlite":
		dbPath := viper.GetString("database.dsn")
		// 确保数据库文件所在的目录存在
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			log.Fatalf("创建数据库目录失败: %v", err)
		}
		DB, err = gorm.Open(sqlite.Open(dbPath), gormConfig)
	case "mysql":
		dsn := viper.GetString("database.dsn")
		DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	default:
		log.Fatalf("不支持的数据库类型: %s", dbType)
	}

	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	log.Println("数据库连接成功!")

	// 自动迁移模式
	runMigrations()

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取数据库实例失败: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

// runMigrations 自动迁移数据库表结构
func runMigrations() {
	log.Println("开始数据库自动迁移...")
	err := DB.AutoMigrate(
		&model.NewsArticle{},
		&model.SiteSettings{},
		&model.Service{},
		&model.CaseStudy{},
		&model.AdminUser{},
		// ... 在这里加入所有定义的模型 ...
	)
	if err != nil {
		log.Fatalf("数据库自动迁移失败: %v", err)
	}
	log.Println("数据库自动迁移成功!")
}
