package ports

import (
	"context"
	"profile-service/internal/domain"

	"github.com/google/uuid"
)

type UserProfileRepository interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile ,error)
	Create(ctx context.Context, userProfile *domain.Profile) error
	Update(ctx context.Context, userProfile *domain.Profile, userID uuid.UUID) error
	Delete(ctx context.Context, userID uuid.UUID) error
}