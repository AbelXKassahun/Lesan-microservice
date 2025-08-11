package app

import (
	"gamification-service/internal/domain"
	"gamification-service/internal/ports"
	"log"
	"time"

	"gorm.io/gorm"
)

type StreakService struct {
	Repo ports.StreakRepository
}

func NewStreakService(repo ports.StreakRepository) *StreakService {
	return &StreakService{Repo: repo}
}

func (s *StreakService) GetStreakByUserID(userID string) (*domain.UserStreak, error) {
	streak, err := s.Repo.GetStreakByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return streak, nil
}

func (s *StreakService) UpdateStreak(userID string) (int, error) {
	today := time.Now().Truncate(24 * time.Hour)

	streak, err := s.Repo.GetStreakByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// first time user
			newStreak := domain.UserStreak{
				UserID:        userID,
				CurrentStreak: 1,
				LastCompleted: today,
			}
			return 1, s.Repo.CreateNewStreak(&newStreak)
		}
		return 0, err
	}

	// Check if already completed today
	if streak.LastCompleted.Equal(today) {
		log.Println("already has a streak")
		return streak.CurrentStreak, nil
	}

	// Check if they were active yesterday
	yesterday := today.AddDate(0, 0, -1)
	if streak.LastCompleted.Equal(yesterday) {
		streak.CurrentStreak += 1
	} else {
		streak.CurrentStreak = 1 // Reset streak
	}
	streak.LastCompleted = today
	
	return streak.CurrentStreak, s.Repo.UpdateStreak(streak)
}