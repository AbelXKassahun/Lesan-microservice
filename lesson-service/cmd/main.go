package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"lesson-service/internal/api"
	"lesson-service/internal/infra/postgres"
	"lesson-service/internal/infra/storage"
	"lesson-service/internal/app"
	"lesson-service/internal/handler"
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
	db := storage.InitPostgres(DB_URL)
	// repos
	userProgressRepo := postgres.NewUserProgressRepository(db)
	unitRepository := postgres.NewUnitRepository(db)
	sectionRepostory := postgres.NewSectionRepository(db)
	lessonRepository := postgres.NewLessonRepository(db)
	exerciseRepo := postgres.NewExerciseRepo(db)
	// services
	exerciseService := app.NewExerciseService(exerciseRepo)
	userProgressService := app.NewUserProgressService(userProgressRepo, unitRepository, sectionRepostory, lessonRepository)
	// handlers
	exerciseHandler := handler.NewExerciseHandler(exerciseService, "exercise-assets")
	userProgressHandler := handler.NewUserProgressHandler(*userProgressService)
	api := api.NewRoutes(exerciseHandler, userProgressHandler)

	return api
}
