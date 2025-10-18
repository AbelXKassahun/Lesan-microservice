package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"lesson-service/internal/app"
	"lesson-service/internal/domain"

	"github.com/google/uuid"
)

type ExerciseHandler struct {
	service     *app.ExerciseService
	assetBucket string
}

func NewExerciseHandler(s *app.ExerciseService, bucket string) *ExerciseHandler {
	return &ExerciseHandler{service: s, assetBucket: bucket}
}

// func (h *ExerciseHandler) RegisterRoutes(mux *http.ServeMux) {
// 	mux.HandleFunc("/exercises/", h.HandleExercise)
// 	mux.HandleFunc("/exercises/lesson/", h.HandleExercisesByLesson)
// }

// dispatch dispatch baby
func (h *ExerciseHandler) HandleExercise(w http.ResponseWriter, r *http.Request) {
	exerciseID := r.URL.Query().Get("exercise_id")
	// parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	switch r.Method {
	case http.MethodGet:
		// exercise id from url
		// if len(parts) < 2 {
		// 	http.Error(w, "missing id", http.StatusBadRequest)
		// 	return
		// }
		// h.GetExercise(w, r, parts[1])

		// exercise id from query param
		if exerciseID == "" {
			http.Error(w, "missing exercise id", http.StatusBadRequest)
			return
		}
		h.getExercise(w, r, exerciseID)
	case http.MethodPost:
		h.createExercise(w, r)
	case http.MethodDelete:
		// exercise id from url
		// if len(parts) < 2 {
		// 	http.Error(w, "missing id", http.StatusBadRequest)
		// 	return
		// }
		// h.DeleteExercise(w, r, parts[1])

		// exercise id from query param
		if exerciseID == "" {
			http.Error(w, "missing exercise id", http.StatusBadRequest)
			return
		}
		h.deleteExercise(w, r, exerciseID)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ExerciseHandler) HandleExercisesByLesson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	// if len(parts) < 3 {
	// 	http.Error(w, "missing lesson_id", http.StatusBadRequest)
	// 	return
	// }
	// h.getExercisesByLesson(w, r, parts[2])

	lesson_id := r.URL.Query().Get("lesson_id")
	if lesson_id == "" {
		http.Error(w, "missing lesson id", http.StatusBadRequest)
		return
	}
	h.getExercisesByLesson(w, r, lesson_id)
}

// handlers
func (h *ExerciseHandler) createExercise(w http.ResponseWriter, r *http.Request) {
	var req domain.Exercise
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exercise, err := h.service.CreateExercise(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, exercise)
}

func (h *ExerciseHandler) getExercise(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	exercise, err := h.service.GetExercise(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Replace asset names with presigned URLs
	h.attachSignedURLs(r.Context(), exercise)

	h.writeJSON(w, exercise)
}

func (h *ExerciseHandler) getExercisesByLesson(w http.ResponseWriter, r *http.Request, lessonIDStr string) {
	lessonID, err := uuid.Parse(lessonIDStr)
	if err != nil {
		http.Error(w, "invalid lesson_id", http.StatusBadRequest)
		return
	}

	exercises, err := h.service.GetExercisesByLesson(r.Context(), lessonID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// For each exercise, replace asset names with presigned URLs
	for i := range *exercises {
		h.attachSignedURLs(r.Context(), &(*exercises)[i])
	}

	h.writeJSON(w, exercises)
}

func (h *ExerciseHandler) deleteExercise(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteExercise(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// helpers

func (h *ExerciseHandler) writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *ExerciseHandler) attachSignedURLs(ctx context.Context, e *domain.Exercise) {
	var data map[string]interface{}
	if err := json.Unmarshal(e.Data, &data); err != nil {
		return
	}

	// key_1 that may contain assets [to be updated]
	key_1 := []string{"prompt_audio_url", "option_audio_url"}
	key_2 := "image_url"
	key_3 := "audio_url"

	for _, key := range key_1 {
		if val, ok := data[key]; ok {
			if assetName, ok := val.(string); ok && !strings.HasPrefix(assetName, "http") {
				// Fetch signed URL
				url, err := app.FetchObjectService(h.assetBucket, assetName)
				if err == nil {
					data[key] = url
				}
			}
		}
	}

	if options, ok := data["options"].([]interface{}); ok {
		for _, option := range options {
			if optMap, ok := option.(map[string]interface{}); ok {
				if assetName, ok := optMap[key_2].(string); ok && !strings.HasPrefix(assetName, "http") {
					// Fetch signed URL
					url, err := app.FetchObjectService(h.assetBucket, assetName)
					if err == nil {
						optMap[key_2] = url
					}
				}
			}
		}
		data["options"] = options
	}

	if items, ok := data["left_items"].([]interface{}); ok {
		firstItem := items[0].(map[string]interface{})
		if _, exists := firstItem["audio_url"]; exists {
			for _, item := range items {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if assetName, ok := itemMap[key_3].(string); ok && !strings.HasPrefix(assetName, "http") {
						// Fetch signed URL
						url, err := app.FetchObjectService(h.assetBucket, assetName)
						if err == nil {
							itemMap[key_2] = url
						}
					}
				}
			}
		}
		data["left_items"] = items
	}

	// Re-marshal back into exercise.Data
	newData, _ := json.Marshal(data)
	e.Data = newData
}
