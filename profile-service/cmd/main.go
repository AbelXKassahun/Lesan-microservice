package main

import (
	"fmt"
	"log"
	"net/http"

	"profile-service/internal/api"
	"profile-service/internal/app"
	"profile-service/internal/handler"
	postgresrepo "profile-service/internal/infra/postgres_repo"
	"profile-service/internal/infra/rabbitmq"
	"profile-service/internal/infra/storage"

	"github.com/joho/godotenv"
)

func main() {
	var myEnv map[string]string
	myEnv, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	port := myEnv["PORT"]

	api, createProfileConsumer, updateProfileConsumer := DependencyInjection(myEnv["DB_URL"])

	RabbitMQHandler(myEnv["RABBITMQ_URL"], createProfileConsumer, updateProfileConsumer)

	log.Printf("Starting the lesan profile service on port %v", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), api.Routes()); err != nil {
		panic(err)
	}
}

func DependencyInjection(DB_URL string) (*api.Routes, *rabbitmq.CreateProfileConsumer, *rabbitmq.UpdateProfileConsumer) {
	db := storage.InitPostgres(DB_URL)

	// repo
	profileRepo := postgresrepo.NewProfileRepository(db)

	//service
	profileService := app.NewProfileService(profileRepo)

	// handler
	profileHandler := handler.NewProfileHandler(profileService)

	// multiplexer
	api := api.NewRoutes(profileHandler)

	//consumers
	createProfileConsumer := rabbitmq.NewCreateProfileConsumer(*profileService)
	updateProfileConsumer := rabbitmq.NewUpdateProfileConsumer(*profileService)

	return api, createProfileConsumer, updateProfileConsumer
}

func RabbitMQHandler(rabbitURL string, createProfileConsumer *rabbitmq.CreateProfileConsumer, updateProfileConsumer *rabbitmq.UpdateProfileConsumer) {
	// rabbitURL := os.Getenv("RABBITMQ_URL")
	err := createProfileConsumer.StartCreateProfileConsumer(rabbitURL)
	if err != nil {
		log.Fatal("❌ Could not start create profile consumer:", err)
	}

	err = updateProfileConsumer.StartUpdateProfileConsumer(rabbitURL)
	if err != nil {
		log.Fatal("❌ Could not start update profile consumer:", err)
	}
}
