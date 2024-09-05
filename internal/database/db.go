package database

import (
	"log"

	"github.com/yanggu0t/go-rdbms-practice/config"
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
