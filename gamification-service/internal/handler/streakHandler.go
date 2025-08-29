package handler

import (
	"encoding/json"
	"net/http"

	"gamification-service/internal/app"
	"gamification-service/internal/utils"
)

type StreakResponse struct {
	UserID string `json:"userId"`
	CurrentStreak int `json:"currentStreak"`
	LastCompleted string `json:"lastCompleted"`
}
type StreakHandler struct {
	StreakService *app.StreakService
}

func NewStreakHandler(streakService *app.StreakService) *StreakHandler {
	return &StreakHandler{StreakService: streakService}
}

func (h *StreakHandler) GetStreakByUserID(w http.ResponseWriter, r *http.Request) {
	var userID string
	userID = r.URL.Query().Get("user_id")
	if userID == "" {
		userID = utils.GetUserFromClaims(w, r)
	}

	streak, err := h.StreakService.GetStreakByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := StreakResponse{
		UserID: streak.UserID,
		CurrentStreak: streak.CurrentStreak,
		LastCompleted: streak.LastCompleted.String(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}