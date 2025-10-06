package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"lesson-service/internal/app"
	"lesson-service/internal/utils"

	"github.com/google/uuid"
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
	user_id := utils.GetUserFromClaims(w, r).NameID

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

func (h *UserProgressHandler) LessonComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req any
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims := utils.GetUserFromClaims(w, r)
	user_id := claims.NameID
	userID, err := uuid.Parse(user_id)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateUserProgress(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// call lesson complete producer below
	// claims.Email

	w.WriteHeader(http.StatusOK)
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
