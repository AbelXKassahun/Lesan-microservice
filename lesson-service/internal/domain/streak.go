package domain

import "time"

// UserStreak represents user's current streak and last active date
type UserStreak struct {
	UserID        string `gorm:"primaryKey;type:text"`
	CurrentStreak int
	LastCompleted time.Time
}
