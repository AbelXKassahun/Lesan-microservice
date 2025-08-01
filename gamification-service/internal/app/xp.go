package app

import (
	"log"

	"gamification-service/internal/ports"
)

type XPService struct {
	Repo ports.XPRepository
}

func NewXPService(repo ports.XPRepository) *XPService {
	return &XPService{Repo: repo}
}

func (s *XPService) HandleLessonCompleted(event LessonCompletedEvent) {
	err := s.Repo.AddXP(event.UserID, event.XP)
	if err != nil {
		log.Printf("❌ Failed to update XP for user %s: %v", event.UserID, err)
		return
	}
	log.Printf("✅ Added %d XP to user %s", event.XP, event.UserID)
}