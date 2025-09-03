package api

import (
	"net/http"
	// "lesson-service/internal/handler"
	"lesson-service/internal/api/middleware"
)

type Routes struct {
	// XPHandler *handler.XPHandler
	// StreakHandler *handler.StreakHandler
	// BadgeHandler *handler.BadgeHandler
	// AggregateHandler *handler.AggregatHandler
}

func NewRoutes() *Routes {
	return &Routes{}
}

func (r *Routes) Routes() *http.ServeMux {
	protected_router := http.NewServeMux()
	main_router := http.NewServeMux()
	main_router.HandleFunc("GET /api/lesson-service/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Gamification Service is up!"))
	})
	// protected_router.HandleFunc("GET /api/lesson-service/lesson", )

	main_router.Handle("/api/lesson-service/", middleware.AuthMiddleware(protected_router))
	return main_router
}