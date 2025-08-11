package handler

import (
	"fmt"
	"gamification-service/internal/app"
	"net/http"
)

type BadgeHandler struct {
	BadgeService *app.BadgeService
}

func NewBadgeHandler(badgeService *app.BadgeService) *BadgeHandler {
	return &BadgeHandler{BadgeService: badgeService}
}

func (h *BadgeHandler) GetBadgesByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	badges, err := h.BadgeService.GetBadgeByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%v", badges))) // might need to be marshaled
}