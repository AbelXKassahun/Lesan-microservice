package postgresrepo

import (
	"context"
	"profile-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileRepo struct {
	DB *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepo {
	return &ProfileRepo{DB: db}
}

func (r *ProfileRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Profile, error) {
	var profile domain.Profile
	// err := r.DB.Where("user_id = ?", userID).First(&profile).Error
	err := r.DB.WithContext(ctx).First(&profile, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepo) Create(ctx context.Context, userProfile *domain.Profile) error {
	return r.DB.WithContext(ctx).Create(userProfile).Error
}

func (r *ProfileRepo) Update(ctx context.Context, userProfile *domain.Profile, userID uuid.UUID) error {
	return r.DB.WithContext(ctx).Where("user_id = ?", userID).Save(userProfile).Error
}

func (r *ProfileRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	// return r.DB.Where("user_id = ?", userID).Delete(&domain.Profile{}).Error
	return r.DB.WithContext(ctx).Delete(&domain.Profile{}, "user_id = ?", userID).Error
}

