package ports

import "gamification-service/internal/domain"

type XPRepository interface {
	GetXPByUserID(userID string) (*domain.UserXP, error)
	AddXP(userID string, xp int) error
}