package postgres

import (
	"gamification-service/internal/domain"

	"gorm.io/gorm"
)

type StreakRepo struct {
	DB *gorm.DB
}

func NewStreakRepo(db *gorm.DB) *StreakRepo {
	return &StreakRepo{DB: db}
}

func (sr *StreakRepo) GetStreakByUserID(userID string) (*domain.UserStreak, error) {
	var streak domain.UserStreak
	result := sr.DB.First(&streak, "user_id = ?", userID)

	return &streak, result.Error
}

func(sr *StreakRepo)CreateNewStreak(newStreak *domain.UserStreak) error {
	return sr.DB.Create(newStreak).Error
}

func(sr *StreakRepo)UpdateStreak(streak *domain.UserStreak) error {
	return sr.DB.Save(streak).Error
}