package handler

import (
	"encoding/json"
	"net/http"

	"profile-service/internal/app"
	"profile-service/internal/domain"
	"profile-service/internal/utils"

	"github.com/google/uuid"
)

type ProfileHandler struct {
	service *app.ProfileService
}

func NewProfileHandler(s *app.ProfileService) *ProfileHandler {
	return &ProfileHandler{service: s}
}


func (h *ProfileHandler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet: 
		h.getProfile(w, r)
	case http.MethodPost:
		h.createProfile(w, r)
	case http.MethodPut:
		h.updateProfile(w, r)
	case http.MethodDelete:
		h.deleteProfile(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
// GET 
func (h *ProfileHandler) getProfile(w http.ResponseWriter, r *http.Request) {
	// ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	// defer cancel()

	user_id := utils.GetUserFromClaims(w, r).NameID
	userID, err := uuid.Parse(user_id)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	profile, err := h.service.GetByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "profile not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

// POST
func (h *ProfileHandler) createProfile(w http.ResponseWriter, r *http.Request) {
	// ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	// defer cancel()

	var req domain.Profile
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// if frontend cant send the userID with the profile objectm then get it from the token yourself
		// user_id := utils.GetUserFromClaims(w, r).NameID
		// userID, err := uuid.Parse(user_id)
		// if err != nil {
		// 	http.Error(w, "invalid user_id", http.StatusBadRequest)
		// 	return
		// }
		// req.UserID = userID

	if err := h.service.CreateProfile(r.Context(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// PUT
func (h *ProfileHandler) updateProfile(w http.ResponseWriter, r *http.Request) {
	// ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	// defer cancel()

	user_id := utils.GetUserFromClaims(w, r).NameID

	userID, err := uuid.Parse(user_id)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	var req domain.Profile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := h.service.UpdateProfile(r.Context(), &req, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DELETE
func (h *ProfileHandler) deleteProfile(w http.ResponseWriter, r *http.Request) {
	// ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	// defer cancel()

	user_id := utils.GetUserFromClaims(w, r).NameID
	userID, err := uuid.Parse(user_id)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteProfile(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
