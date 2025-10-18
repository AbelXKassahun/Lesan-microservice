package storage

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"profile-service/internal/domain"
)

func InitPostgres(connectionString string) *gorm.DB {
	log.Println(connectionString)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ failed to connect to DB: %v", err)
	}

	if err := domain.Migrate(db); err != nil {
		log.Fatalf("❌ migration failed: %v", err)
	}
	log.Println("✅ DB migrated successfully")

	return db
}
