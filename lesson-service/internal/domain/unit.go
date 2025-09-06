package domain

import (
	"time"

	"github.com/google/uuid"
)

type Unit struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title        string    `gorm:"type:text;not null" json:"title"`
	Description  string    `gorm:"type:text" json:"description"`
	SectionCount int       `gorm:"type:int" json:"section_count"`
	OrderIndex   int       `gorm:"type:int;not null" json:"order_index"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
