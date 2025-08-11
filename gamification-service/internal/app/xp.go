package app

import (
	"log"

	"gamification-service/internal/domain"
	"gamification-service/internal/ports"
)

type XPService struct {
	Repo ports.XPRepository
}

func NewXPService(repo ports.XPRepository) *XPService {
	return &XPService{Repo: repo}
}

func (s *XPService) UpdateXP(userID string, xp int) {
	// update XP on lesson complete`
	// err := s.Repo.AddXP(event.UserID, event.XP)`
	err := s.Repo.AddXP(userID, xp)
	if err != nil {
		log.Printf("❌ Failed to update XP for user %s: %v", userID, err)
		return
	}
	log.Printf("✅ Added %d XP to user %s", xp, userID)

	//
}

func (s *XPService) GetXPByUserID(userID string) (*domain.UserXP, error) {
	return s.Repo.GetXPByUserID(userID)
}
func (s *XPService) GetUsersByLeague(league string) (*[]domain.UserXP, error) {
	return s.Repo.GetUsersByLeague(league)
}