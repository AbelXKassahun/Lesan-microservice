package handler

import (
	"encoding/json"
	"net/http"

	"gamification-service/internal/app"
	"gamification-service/internal/utils"
)

type XPResponse struct {
	UserID string `json:"userID"`
	TotalXP int `json:"totalXP"`
	League string `json:"league"`
}

type XPHandler struct {
	XPService *app.XPService
}

func NewXPHandler(xpService *app.XPService) *XPHandler {
	return &XPHandler{XPService: xpService}
}

func (h *XPHandler) GetXPByUserID(w http.ResponseWriter, r *http.Request) {
	var userID string
	userID = r.URL.Query().Get("user_id")
	if userID == "" {
		userID = utils.GetUserFromClaims(w, r)
	}

	xp, err := h.XPService.GetXPByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := XPResponse {
		UserID: xp.UserID,
		TotalXP: xp.Total,
		League: string(xp.League),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *XPHandler) GetUsersByLeague(w http.ResponseWriter, r *http.Request) {
	league := r.URL.Query().Get("league")
	xp, err := h.XPService.GetUsersByLeague(league)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response  []XPResponse
	for _, row := range *xp {
		response = append(response, XPResponse{
				UserID: row.UserID,
				TotalXP: row.Total,
				League: string(row.League),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
