package config

import (
	"chatroom/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(configPath string) {
	// 调用 LoadConfig 时处理 error
	config, err := LoadConfig(configPath)
	log.Println("解析成功")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 使用 GORM 连接 MySQL
	db, err := gorm.Open(mysql.Open(config.Database.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	DB = db
	log.Println("Database connected successfully!")

	// 自动迁移表结构
	err = DB.AutoMigrate(
		&models.AccountInfo{},
		&models.Contacts{},
		&models.SystemSetting{},
		&models.ApplyInfo{},
		&models.GroupChatInfo{},
		&models.GroupMemberInfo{},
		&models.Notes{},
		&models.Favorites{},
		&models.Codes{},
		&models.DDLS{},
		&models.Recycle{},
		&models.GroupDivide{},
		&models.FriendDivide{},
		&models.NoteDivide{},
		&models.Links{},
		&models.MessageInfo{},
		&models.DeleteInfo{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database tables: %v", err)
	}
}

// GetDB 返回数据库实例
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database connection is not initialized. Call InitDB first.")
	}
	return DB
}
