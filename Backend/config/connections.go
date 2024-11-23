package config

import (
	"chatroom/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	// 数据库连接配置
	dsn := "user:password@tcp(127.0.0.1:3306)/chatroom?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
	log.Println("Database tables migrated successfully!")
}

// GetDB 返回数据库实例
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database connection is not initialized. Call InitDB first.")
	}
	return DB
}
