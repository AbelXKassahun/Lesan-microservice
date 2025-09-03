package main

import (
	"fmt"
	"log"
	"net/http"

	"lesson-service/internal/api"

	"github.com/joho/godotenv"
	// "lesson-service/internal/app"
	// "lesson-service/internal/handler"
	// "lesson-service/internal/infra/postgres"
)

func main() {
	var myEnv map[string]string
	myEnv, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	port := myEnv["PORT"]

	api := DependencyInjection(myEnv["DB_URL"])

	log.Printf("Starting the lesan lesson service on port %v", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), api.Routes()); err != nil {
		panic(err)
	}
}

func DependencyInjection(DB_URL string) *api.Routes {
	// db := storage.InitPostgres(DB_URL)
	// // repos
	// xpRepo := postgres.NewXPRepo(db)
	// streakRepo := postgres.NewStreakRepo(db)
	// badgeRepo := postgres.NewBadgeRepo(db)
	// // services
	// xpService := app.NewXPService(xpRepo)
	// streakService := app.NewStreakService(streakRepo)
	// badgeService := app.NewBadgeService(badgeRepo)
	// // handlers
	// xpHandler := handler.NewXPHandler(xpService)
	// streakHandler := handler.NewStreakHandler(streakService)
	// badgeHandler := handler.NewBadgeHandler(badgeService)
	// aggregateHandler := handler.NewAggregateHandler(xpService, streakService, badgeService)

	api := api.NewRoutes()

	return api
}
