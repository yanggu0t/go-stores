package database

import (
	"log"

	"github.com/yanggu0t/go-rdbms-practice/config"
	"github.com/yanggu0t/go-rdbms-practice/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		PrepareStmt: false, // 禁用預處理語句緩存
	})
	if err != nil {
		log.Fatalf("無法連接到數據庫：%v", err)
	}

	log.Println("數據庫初始化成功")
	return db
}

func AutoMigrate(db *gorm.DB) {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Role{},
		&models.Permission{},
		&models.ProjectUserRole{},
	)
	if err != nil {
		log.Fatalf("Error during database migration: %v", err)
	}

	log.Println("Database migration completed successfully")
}
