package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"lesson-service/internal/domain"
)

type UserProgressRepository struct {
	DB *gorm.DB
}

func NewUserProgressRepository(db *gorm.DB) *UserProgressRepository {
	return &UserProgressRepository{DB: db}
}

func (r *UserProgressRepository) Create(ctx context.Context, userProgress *domain.UserProgress) error {
	return r.DB.WithContext(ctx).Create(userProgress).Error
}

func (r *UserProgressRepository) GetByUser(ctx context.Context, userID uuid.UUID) (*domain.UserProgress, error) {
	var up domain.UserProgress
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).First(&up).Error
	if err != nil {
		return nil, err
	}
	return &up, nil
}

func (r *UserProgressRepository) Update(ctx context.Context, userProgress *domain.UserProgress) error {
	return r.DB.WithContext(ctx).Save(userProgress).Error
}
