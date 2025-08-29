package handler

import (
	"fmt"
	"net/http"

	"lesson-service/internal/app"
)

type XPHandler struct {
	XPService *app.XPService
}

func NewXPHandler(xpService *app.XPService) *XPHandler {
	return &XPHandler{XPService: xpService}
}

func (h *XPHandler) GetXPByUserID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	xp, err := h.XPService.GetXPByUserID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%d", xp.Total)))
}

func (h *XPHandler) GetUsersByLeague(w http.ResponseWriter, r *http.Request) {
	league := r.URL.Query().Get("league")
	xp, err := h.XPService.GetUsersByLeague(league)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%v", xp)))
}
