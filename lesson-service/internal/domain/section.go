package domain

import (
	"time"

	"github.com/google/uuid"
)

type Section struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UnitID      uuid.UUID `gorm:"type:uuid;not null" json:"unit_id"`
	Type        string    `gorm:"type:text" json:"type"`
	LessonCount int       `gorm:"type:int" json:"lesson_count"`
	OrderIndex  int       `gorm:"type:int" json:"order_index"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
