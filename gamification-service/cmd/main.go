package main

import (
		"flag"
		"fmt"
		"log"
		"net/http"

	"gamification-service/internal/api"
	"gamification-service/internal/app"
	"gamification-service/internal/handler"
	"gamification-service/internal/infra/postgres"
	"gamification-service/internal/infra/rabbitmq"
	"gamification-service/internal/infra/storage"

	"github.com/joho/godotenv"
)

func main() {
	var myEnv map[string]string
	myEnv, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	port := initializePortFlag(myEnv["PORT"])

	api, lessonCompletedConsumer := DependencyInjection(myEnv["DB_URL"])

	RabbitMQHandler(myEnv["RABBITMQ_URL"], lessonCompletedConsumer)

	// http.Handle("GET /metadata", http.HandlerFunc(h.GetMetadata))

	log.Printf("Starting the lesan gamification service on port %v", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), api.Routes()); err != nil {
		panic(err)
	}
}

func DependencyInjection(DB_URL string) (*api.RoutesType, *rabbitmq.LessonCompletedConsumer) {
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
	lessonCompletedConsumer := rabbitmq.NewLessonCompletedConsumer(xpService, streakService, badgeService)

	return api, lessonCompletedConsumer
}

func initializePortFlag(defaultPort string) string {
	var port string
	flag.StringVar(&port, "port", defaultPort, "API handler port")
	flag.Parse()
	return port
}

func RabbitMQHandler(rabbitURL string, lessonCompletedConsumer *rabbitmq.LessonCompletedConsumer) {
	// rabbitURL := os.Getenv("RABBITMQ_URL")
	err := lessonCompletedConsumer.StartLessonCompletedConsumer(rabbitURL)
	if err != nil {
		log.Fatal("❌ Could not start RabbitMQ consumer:", err)
	}
}