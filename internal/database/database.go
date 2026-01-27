package database

import (
	"errors"
	"log"

	"luggage-sys2/internal/config"
	"luggage-sys2/internal/models"
	"luggage-sys2/internal/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(mysql.Open(config.DBDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 修复取件码索引：先删除唯一索引（如果存在），再执行迁移
	fixRetrievalCodeIndex()

	// 自动迁移
	err = DB.AutoMigrate(
		&models.User{},
		&models.Luggage{},
		&models.Storeroom{},
		&models.StoredLog{},
		&models.UpdatedLog{},
		&models.RetrievedLog{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 初始化默认数据
	initDefaultData()
}

func initDefaultData() {
	// 创建默认管理员（若不存在）
	var admin models.User
	err := DB.Where("username = ?", "admin").First(&admin).Error
	if err == nil {
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("Failed to query default admin user:", err)
		return
	}

	hashedPassword, _ := utils.HashPassword("123456")
	defaultUser := models.User{
		Username: "admin",
		Password: hashedPassword,
		Role:     "admin",
		HotelID:  1,
	}
	if err := DB.Create(&defaultUser).Error; err != nil {
		log.Println("Failed to create default admin user:", err)
		return
	}
	log.Println("Default user created: admin / 123456")
}

// fixRetrievalCodeIndex 修复取件码索引：从唯一索引改为普通索引
func fixRetrievalCodeIndex() {
	// 检查表是否存在
	var tableExists int
	DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'luggages'").Scan(&tableExists)
	if tableExists == 0 {
		return // 表不存在，AutoMigrate 会创建
	}

	// 查询所有可能的唯一索引名称
	var indexes []struct {
		KeyName   string `gorm:"column:Key_name"`
		NonUnique int    `gorm:"column:Non_unique"`
	}
	
	err := DB.Raw(`
		SELECT Key_name, Non_unique 
		FROM information_schema.STATISTICS 
		WHERE table_schema = DATABASE() 
		AND table_name = 'luggages' 
		AND column_name = 'retrieval_code'
		AND Non_unique = 0
	`).Scan(&indexes).Error
	
	if err == nil && len(indexes) > 0 {
		// 删除所有唯一索引
		for _, idx := range indexes {
			log.Printf("Removing unique index '%s' on retrieval_code...", idx.KeyName)
			if err := DB.Exec("ALTER TABLE luggages DROP INDEX `" + idx.KeyName + "`").Error; err != nil {
				log.Printf("Warning: Failed to drop index %s: %v", idx.KeyName, err)
			} else {
				log.Printf("Successfully removed unique index '%s'", idx.KeyName)
			}
		}
	}
	
	// AutoMigrate 会自动创建普通索引（因为模型定义中已经是 index 而不是 uniqueIndex）
}
