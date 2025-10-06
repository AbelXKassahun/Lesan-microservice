package app

import (
	"gamification-service/internal/domain"
	"gamification-service/internal/ports"
	"log"
	"time"

	"gorm.io/gorm"
)

type BadgeService struct {
	Repo ports.BadgeRepository
}
type ReadableBadge struct {
	ID        uint
	UserID    string
	BadgeName string
	AwardedAt time.Time
}

func NewBadgeService(repo ports.BadgeRepository) *BadgeService {
	return &BadgeService{Repo: repo}
}

func (s *BadgeService) GetBadgeByUserID(userID string) (*[]domain.UserBadge, error) {
	badges, err := s.Repo.GetBadgesByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return badges, nil
}

func (s *BadgeService) TryAwardBadge(userID string, badgeName string) (bool, error) {
	err := s.Repo.CheckIfBadgeAwarded(userID, badgeName)
	badgeAwarded := false
	if err == nil {
		// Already has this badge
		log.Printf("already has this badge, %v", badgeName)
		return badgeAwarded, nil
	}
	if err != gorm.ErrRecordNotFound {
		return badgeAwarded, err
	}

	badgeAwarded = true
	newBadge := domain.UserBadge{
		UserID:    userID,
		BadgeName: badgeName,
		AwardedAt: time.Now(),
	}
	return badgeAwarded, s.Repo.AddBadge(&newBadge)
}
