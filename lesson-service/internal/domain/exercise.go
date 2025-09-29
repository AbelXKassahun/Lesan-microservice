package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Exercise struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	LessonID    uuid.UUID      `gorm:"type:uuid;not null" json:"lesson_id"`
	Type        string         `gorm:"type:text;not null" json:"type"`
	Subtype     *string        `gorm:"type:text" json:"subtype,omitempty"`
	Instruction *string        `gorm:"type:text" json:"instruction,omitempty"`
	Data        datatypes.JSON `gorm:"type:jsonb;not null" json:"data"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
}
