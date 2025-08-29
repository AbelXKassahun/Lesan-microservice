package storage

import (
	"log"
	// "os"

	// "github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "github.com/joho/godotenv"

	"lesson-service/internal/domain"
)

func InitPostgres(connectionString string) *gorm.DB { // *pgx.Conn
	// connStr := "postgres://user:password@postgres:5432/finTechDB?sslmode=disable"
	// connStr := strings.TrimSpace(os.Getenv("DB_URL"))

	// godotenv.Load("../../../.env")
	// db, err := pgx.Connect(context.Background(), os.Getenv(connectionString))
	// if err != nil {
	// 	log.Fatal("❌ Could not connect to Postgres:", err)
	// }

	// dsn := os.Getenv("DB_URL")
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
