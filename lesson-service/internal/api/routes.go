package api

import (
	"net/http"
	"lesson-service/internal/api/middleware"
	"lesson-service/internal/handler"
)

type Routes struct {
	ExerciseHandler *handler.ExerciseHandler
	UserProgressHandler *handler.UserProgressHandler
}

func NewRoutes(
	exerciseHandler *handler.ExerciseHandler, 
	userProgressHandler *handler.UserProgressHandler,
) *Routes {
	return &Routes{
		ExerciseHandler: exerciseHandler,
		UserProgressHandler: userProgressHandler,
	}
}


func (r *Routes) Routes() *http.ServeMux {
	protected_router := http.NewServeMux()
	main_router := http.NewServeMux()
	main_router.HandleFunc("GET /api/lesson-service/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Lesson Service is up!"))
	})

	protected_router.HandleFunc("/api/lesson-service/exercises/", r.ExerciseHandler.HandleExercise)
	
	protected_router.HandleFunc("/api/lesson-service/exercises/lesson", r.ExerciseHandler.HandleExercisesByLesson)

	protected_router.HandleFunc("/api/lesson-service/user-progress", r.UserProgressHandler.GetUserProgress)
	
	main_router.Handle("/api/lesson-service/", middleware.AuthMiddleware(protected_router))
	return main_router
}