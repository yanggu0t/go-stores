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

	// Disable foreign key checks during migration
	db.Exec("SET CONSTRAINTS ALL DEFERRED;")

	// Define the order of migrations
	migrations := []struct {
		Model interface{}
		Name  string
	}{
		{&models.User{}, "User"},
		{&models.Permission{}, "Permission"},
		{&models.Role{}, "Role"},
		{&models.Project{}, "Project"},
		{&models.UserSession{}, "UserSession"},
		{&models.ProjectUserRole{}, "ProjectUserRole"},
		{&models.ProjectPermission{}, "ProjectPermission"},
		{&models.RolePermission{}, "RolePermission"},
	}

	// Perform migrations in order
	for _, migration := range migrations {
		log.Printf("Migrating %s...", migration.Name)
		if err := db.AutoMigrate(migration.Model); err != nil {
			log.Fatalf("Error migrating %s: %v", migration.Name, err)
		}
	}

	// Re-enable foreign key checks
	db.Exec("SET CONSTRAINTS ALL IMMEDIATE;")

	log.Println("Database migration completed successfully")
}
