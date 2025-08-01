package storage

import (
	"context"
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5"
)


func InitPostgres(connectionString string) *pgx.Conn {
	// connStr := "postgres://user:password@postgres:5432/finTechDB?sslmode=disable"
	// connStr := strings.TrimSpace(os.Getenv("DB_URL"))
	
	godotenv.Load()
	db, err := pgx.Connect(context.Background(), os.Getenv(connectionString))
	if err != nil {
		log.Fatal("❌ Could not connect to Postgres:", err)
	}
	return db
}
