package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"gorm.io/gorm"

	"gamification-service/internal/app"
	"gamification-service/internal/domain"
	"gamification-service/internal/utils"
)

type AggregatHandler struct {
	XPService     *app.XPService
	StreakService *app.StreakService
	BadgeService  *app.BadgeService
}

type Stats struct {
	XP     *domain.UserXP      `json:"xp"`
	Streak *domain.UserStreak  `json:"streak"`
	Badge  *[]domain.UserBadge `json:"badge"`
}

func NewAggregateHandler(xpService *app.XPService,
	streakService *app.StreakService,
	badgeService *app.BadgeService) *AggregatHandler {
	return &AggregatHandler{
		XPService:     xpService,
		StreakService: streakService,
		BadgeService:  badgeService,
	}
}

func (h *AggregatHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	var response Stats
	var err error
	var userID string

	userID = r.URL.Query().Get("user_id")
	if userID == "" {
		userID = utils.GetClaimsFromToken(w, r).NameID
	}

	xp, err := h.XPService.GetXPByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("xp")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	streak, err := h.StreakService.GetStreakByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("streak")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	badges, err := h.BadgeService.GetBadgeByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("badges")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response = Stats{
		XP:     xp,
		Streak: streak,
		Badge:  badges,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
