package api

import (
	"net/http"

	"profile-service/internal/api/middleware"
	"profile-service/internal/handler"
)

type Routes struct {
	ProfileHandler *handler.ProfileHandler
}

func NewRoutes(
	profileHandler *handler.ProfileHandler,
) *Routes {
	return &Routes{
		ProfileHandler: profileHandler,
	}
}

func (r *Routes) Routes() *http.ServeMux {
	protected_router := http.NewServeMux()
	main_router := http.NewServeMux()
	main_router.HandleFunc("GET /api/profile-service/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Profile Service is up!"))
	})

	protected_router.HandleFunc("/api/profile-service/user-profile", r.ProfileHandler.HandleProfile)

	main_router.Handle("/api/profile-service/", middleware.AuthMiddleware(protected_router))
	return main_router
}
