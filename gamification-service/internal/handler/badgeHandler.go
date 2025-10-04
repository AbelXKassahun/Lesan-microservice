package handler

import (
	"encoding/json"
	"net/http"

	"gamification-service/internal/app"
	"gamification-service/internal/utils"
)

type BadgeResponse struct {
	BadgeID   int    `json:"badge_id"`
	UserID    string `json:"user_id"`
	BadgeName string `json:"badge_name"`
	AwardedAt string `json:"awarded_at"`
}
type BadgeHandler struct {
	BadgeService *app.BadgeService
}

func NewBadgeHandler(badgeService *app.BadgeService) *BadgeHandler {
	return &BadgeHandler{BadgeService: badgeService}
}

func (h *BadgeHandler) GetBadgesByUserID(w http.ResponseWriter, r *http.Request) {
	var userID string
	userID = r.URL.Query().Get("user_id")
	if userID == "" {
		userID = utils.GetClaimsFromToken(w, r).NameID
	}

	badges, err := h.BadgeService.GetBadgeByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response []BadgeResponse
	for _, row := range *badges {
		response = append(response, BadgeResponse{
			BadgeID:   int(row.ID),
			UserID:    row.UserID,
			BadgeName: row.BadgeName,
			AwardedAt: row.AwardedAt.String(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
