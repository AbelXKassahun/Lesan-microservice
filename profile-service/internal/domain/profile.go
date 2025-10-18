package domain

import (
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Profile struct {
	ID               uint      `gorm:"primaryKey;autoIncrement"`
	UserID           uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"` // FK to Identity User
	FullName         string    `gorm:"type:varchar(150)"`
	UserName         string    `gorm:"type:varchar(150)"`
	Email            string    `gorm:"type:varchar(150)"`
	DaysActive       int       `gorm:"default:0"`
	Streak           int       `gorm:"default:0"`
	XP               int       `gorm:"default:0"`
	LessonsCompleted int       `gorm:"default:0"`
	MemberSince      time.Time `gorm:"autoCreateTime"`
	League           string    `gorm:"type:varchar(50)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
