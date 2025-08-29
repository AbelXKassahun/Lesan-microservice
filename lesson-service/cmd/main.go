package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"lesson-service/internal/api"
	"lesson-service/internal/app"
	"lesson-service/internal/handler"
	"lesson-service/internal/infra/postgres"
	"lesson-service/internal/infra/storage"
)

func main () {
	var myEnv map[string]string
	myEnv, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	
	port := initializePortFlag(myEnv["PORT"])

	api := DependencyInjection(myEnv["DB_URL"])

	log.Printf("Starting the lesan gamification service on port %v", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), api.Routes()); err != nil {
		panic(err)
	}
}

func initializePortFlag(defaultPort string) string {
	var port string
	flag.StringVar(&port, "port", defaultPort, "API handler port")
	flag.Parse()
	return port
}

func DependencyInjection(DB_URL string) (*api.RoutesType) {
	db := storage.InitPostgres(DB_URL)
	// repos
	xpRepo := postgres.NewXPRepo(db)
	streakRepo := postgres.NewStreakRepo(db)
	badgeRepo := postgres.NewBadgeRepo(db)
	// services 
	xpService := app.NewXPService(xpRepo)
	streakService := app.NewStreakService(streakRepo)
	badgeService := app.NewBadgeService(badgeRepo)
	// handlers
	xpHandler := handler.NewXPHandler(xpService)
	streakHandler := handler.NewStreakHandler(streakService)
	badgeHandler := handler.NewBadgeHandler(badgeService)
	aggregateHandler := handler.NewAggregateHandler(xpService, streakService, badgeService)
	
	api := api.NewRoutes(xpHandler, streakHandler, badgeHandler, aggregateHandler)

	return api
}