package handler

import "net/http"

type LessonHandler struct {

}

func NewLessonHandler() *LessonHandler{
	return &LessonHandler{}
}

func (h *LessonHandler) GetNextLesson(w http.ResponseWriter, r *http.Request) {

}