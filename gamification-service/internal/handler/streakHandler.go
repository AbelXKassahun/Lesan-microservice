package handler

import (
	"fmt"
	"gamification-service/internal/app"
	"net/http"
)

type StreakHandler struct {
	StreakService *app.StreakService
}

func NewStreakHandler(streakService *app.StreakService) *StreakHandler {
	return &StreakHandler{StreakService: streakService}
}

func (h *StreakHandler) GetStreakByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	streak, err := h.StreakService.GetStreakByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%v", streak))) // might need to be marshaled
}