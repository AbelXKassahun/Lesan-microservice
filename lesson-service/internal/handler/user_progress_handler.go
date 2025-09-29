package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"lesson-service/internal/app"
	"lesson-service/internal/utils"
)

type UserProgressHandler struct {
	service app.UserProgressService
}

func NewUserProgressHandler(s app.UserProgressService) *UserProgressHandler {
	return &UserProgressHandler{service: s}
}

func (h *UserProgressHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/user-progress/", h.GetUserProgress)
}

func (h *UserProgressHandler) GetUserProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Expect /user-progress/{user_id}
	// parts := splitPath(r.URL.Path)
	// if len(parts) < 2 {
	// 	http.Error(w, "missing user_id", http.StatusBadRequest)
	// 	return
	// }
	user_id := utils.GetUserFromClaims(w, r)
	
	userID, err := uuid.Parse(user_id)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	result, err := h.service.GetUserProgress(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// helper
func splitPath(p string) []string {
	out := []string{}
	for _, part := range strings.Split(strings.Trim(p, "/"), "/") {
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
