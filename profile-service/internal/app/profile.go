package app

import (
	"context"
	"log"

	"github.com/google/uuid"

	"profile-service/internal/domain"
	"profile-service/internal/ports"
)

type ProfileService struct {
	Repo ports.UserProfileRepository
}

func NewProfileService(repo ports.UserProfileRepository) *ProfileService {
	return &ProfileService{Repo: repo}
}

func (s *ProfileService) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	return s.Repo.GetByUserID(ctx, userID)
}

func (s *ProfileService) CreateProfile(ctx context.Context, profile *domain.Profile) error {
	// optional: check if already exists
	existing, err := s.Repo.GetByUserID(ctx, profile.UserID)
	if err == nil && existing != nil {
		log.Println("profile already exists")
		return nil
	}
	return s.Repo.Create(ctx, profile)
}

func (s *ProfileService) UpdateProfile(ctx context.Context, profile *domain.Profile, userID uuid.UUID) error {
	return s.Repo.Update(ctx, profile, userID)
}

func (s *ProfileService) DeleteProfile(ctx context.Context, userID uuid.UUID) error {
	return s.Repo.Delete(ctx, userID)
}
