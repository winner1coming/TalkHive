package config

import (
	"TalkHive/global"
	"TalkHive/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	// 确保配置文件中的 DSN 已正确加载
	dsn := AppConfig.Database.Dsn
	if dsn == "" {
		log.Fatal("DSN为空")
	}

	// 使用 GORM 连接 MySQL
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	log.Println("数据库连接成功")

	// 配置数据库连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("获取数据库连接池失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns) // 从 AppConfig 加载最大空闲连接数
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns) // 从 AppConfig 加载最大打开连接数

	// 自动迁移表结构
	err = autoMigrateTables()
	if err != nil {
		log.Fatalf("自动迁移表失败: %v", err)
	}
	global.Db = DB
}

// autoMigrateTables 自动迁移数据库表结构
func autoMigrateTables() error {
	// 自动迁移的表列表
	tables := []interface{}{
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
	}

	// 执行自动迁移
	for _, table := range tables {
		err := DB.AutoMigrate(table)
		if err != nil {
			return err
		}
	}
	log.Println("表迁移成功")
	return nil
}

// GetDB 返回全局数据库实例
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("数据库实例为空")
	}
	return DB
}
