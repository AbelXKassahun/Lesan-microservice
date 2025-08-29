package api

import (
	"gamification-service/internal/handler"
	"gamification-service/internal/api/middleware"
	"net/http"
)

type RoutesType struct {
	XPHandler *handler.XPHandler
	StreakHandler *handler.StreakHandler
	BadgeHandler *handler.BadgeHandler
	AggregateHandler *handler.AggregatHandler
}

func NewRoutes(xpHandler *handler.XPHandler, streakHandler *handler.StreakHandler, badgeHandler *handler.BadgeHandler, aggregateHandler *handler.AggregatHandler) *RoutesType {
	return &RoutesType{
		XPHandler: xpHandler,
		StreakHandler: streakHandler,
		BadgeHandler: badgeHandler,
		AggregateHandler: aggregateHandler,
	}
}

func (r *RoutesType) Routes() *http.ServeMux {
	mainRouter := http.NewServeMux() 
	protected_router := http.NewServeMux()
	
	mainRouter.HandleFunc("GET /api/game-service/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Gamification Service is up!"))
	})
	protected_router.HandleFunc("GET /api/game-service/xp", r.XPHandler.GetXPByUserID)
	protected_router.HandleFunc("GET /api/game-service/streak", r.StreakHandler.GetStreakByUserID)
	protected_router.HandleFunc("GET /api/game-service/badge", r.BadgeHandler.GetBadgesByUserID)

	// most likely to be used a lot
	protected_router.HandleFunc("GET /api/game-service/stats", r.AggregateHandler.GetUserStats)
	protected_router.HandleFunc("GET /api/game-service/league", r.XPHandler.GetUsersByLeague)

	mainRouter.Handle("/api/game-service/", middleware.AuthMiddleware(protected_router))
	return mainRouter
}