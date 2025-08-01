package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"gamification-service/internal/api"
	"gamification-service/internal/app"
	"gamification-service/internal/infra/postgres"
	"gamification-service/internal/infra/rabbitmq"
	"gamification-service/internal/infra/storage"

	"github.com/joho/godotenv"
)

func main() {
	var myEnv map[string]string
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := initializePortFlag(myEnv["PORT"])

	db := storage.InitPostgres(myEnv["DB_URL"])
	xpRepo := postgres.NewXPRepo(db)
	xpService := app.NewXPService(xpRepo)

	lessonCompleteRabbitMQHandler(myEnv["RABBITMQ_URL"], xpService)

	// http.Handle("GET /metadata", http.HandlerFunc(h.GetMetadata))

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

func lessonCompleteRabbitMQHandler(rabbitURL string, xpService *app.XPService) {
	// rabbitURL := os.Getenv("RABBITMQ_URL")
	err := rabbitmq.StartLessonCompletedConsumer(rabbitURL, xpService)
	if err != nil {
		log.Fatal("❌ Could not start RabbitMQ consumer:", err)
	}
}
