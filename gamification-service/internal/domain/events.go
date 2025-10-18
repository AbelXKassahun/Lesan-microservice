package domain

type LessonCompletedEvent struct {
	Event     string `json:"event"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	LessonID  string `json:"lesson_id"`
	Timestamp string `json:"timestamp"`
	XP        int    `json:"xp"`
}
