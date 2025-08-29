package ports

import "lesson-service/internal/domain"

type XPRepository interface {
	GetXPByUserID(userID string) (*domain.UserXP, error)
	AddXP(userID string, xp int) error
	GetUsersByLeague(league string) (*[]domain.UserXP, error)
}

type StreakRepository interface {
	GetStreakByUserID(userID string) (*domain.UserStreak, error)
	CreateNewStreak(newStreak *domain.UserStreak) error
	UpdateStreak(streak *domain.UserStreak) error
}

type BadgeRepository interface {
	GetBadgesByUserID(userID string) (*[]domain.UserBadge, error)
	CheckIfBadgeAwarded(userID string, badgeName string)  error
	AddBadge(badge *domain.UserBadge) error
}