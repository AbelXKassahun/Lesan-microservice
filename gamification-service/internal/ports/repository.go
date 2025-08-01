package ports

import "gamification-service/internal/models"

type XPRepository interface {
	GetXPByUserID(userID string) (*models.UserXP, error)
	AddXP(userID string, xp int) error
}