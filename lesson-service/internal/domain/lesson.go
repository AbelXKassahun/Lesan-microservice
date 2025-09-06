package domain

import (
	"time"

	"github.com/google/uuid"
)

type Lesson struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UnitID        uuid.UUID `gorm:"type:uuid;not null" json:"unit_id"`
	SectionID     uuid.UUID `gorm:"type:uuid;not null" json:"section_id"`
	Title         string    `gorm:"type:text" json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	ExerciseCount int       `gorm:"type:int" json:"exercise_count"`
	OrderIndex    int       `gorm:"type:int" json:"order_index"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}
