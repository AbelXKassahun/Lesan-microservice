package api

import (
	"lesson-service/internal/handler"
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
	router := http.NewServeMux()
	router.HandleFunc("/api/game-service/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Gamification Service is up!"))
	})
	router.HandleFunc("GET /api/game-service/xp", r.XPHandler.GetXPByUserID)
	router.HandleFunc("GET /api/game-service/streak", r.StreakHandler.GetStreakByUserID)
	router.HandleFunc("GET /api/game-service/badge", r.BadgeHandler.GetBadgesByUserID)

	// most likely to be used a lot
	router.HandleFunc("GET /api/game-service/stats", r.AggregateHandler.GetUserStats)
	router.HandleFunc("GET /api/game-service/league", r.XPHandler.GetUsersByLeague)
	
	return router
}