package domain

import (
	// "gorm.io/gorm"
	"time"
)

// UserXP represents total XP a user has earned
type UserXP struct {
	UserID string `gorm:"primaryKey;type:text"`
	Total  int
}

// UserStreak represents user's current streak and last active date
type UserStreak struct {
	UserID      string    `gorm:"primaryKey;type:text"`
	CurrentStreak int
	LastActive  time.Time
}

// UserBadge represents badges earned by a user
type UserBadge struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    string    `gorm:"index;type:text"`
	BadgeName string
	AwardedAt time.Time
}