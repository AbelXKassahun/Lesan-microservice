package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserProgress struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	CurrentUnit   uuid.UUID `gorm:"type:uuid;not null" json:"current_unit"`
	CurrentSection uuid.UUID `gorm:"type:uuid;not null" json:"current_section"`
	CurrentLesson uuid.UUID `gorm:"type:uuid;not null" json:"current_lesson"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// this fixes the "user_progresses does not exist" issue
func (UserProgress) TableName() string {
	return "user_progress"
}