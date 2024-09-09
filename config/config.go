package config

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	JWTSecret   string
}

func Load() *Config {
	// 獲取當前文件的路徑
	_, filename, _, _ := runtime.Caller(0)
	// 獲取項目根目錄的路徑
	rootDir := filepath.Join(filepath.Dir(filename), "..")

	// 加載 .env 文件，指定完整路徑
	err := godotenv.Load(filepath.Join(rootDir, ".env"))
	if err != nil {
		log.Println("Warning: Error loading .env file")
	}

	dbURL := getEnv("DATABASE_URL", "")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	// 解析 URL
	u, err := url.Parse(dbURL)
	if err != nil {
		log.Fatalf("Invalid database URL: %v", err)
	}

	// 從 URL 中提取密碼
	password, _ := u.User.Password()

	// 重新構建連接字符串
	connStr := "postgres://" + u.User.Username() + ":" + password + "@" + u.Host + u.Path

	return &Config{
		DatabaseURL: connStr,
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
