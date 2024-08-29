package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// 獲取當前工作目錄
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("無法獲取當前工作目錄：%v", err)
	}

	// 構建 .env 文件的完整路徑
	envPath := filepath.Join(currentDir, ".env")

	// 加載 .env 文件
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("警告: 無法加載 .env 文件（%s）。使用環境變量。錯誤：%v", envPath, err)
	}

	// 獲取數據庫 URL
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("錯誤：DATABASE_URL 環境變量未設置")
	}

	// 初始化 GORM
	var dbErr error
	DB, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false, // 禁用預處理語句緩存
	})

	if dbErr != nil {
		log.Fatalf("無法連接到數據庫：%v", dbErr)
	}

	log.Println("數據庫初始化成功")
}
