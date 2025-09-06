package postgres

import (
	"context"
	"lesson-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExerciseRepo struct {
	DB *gorm.DB
}

func NewExerciseRepo (db *gorm.DB) *ExerciseRepo {
	return &ExerciseRepo{DB: db}
}

func (r *ExerciseRepo) Create(ctx context.Context, exercise *domain.Exercise) error {
	return r.DB.WithContext(ctx).Create(exercise).Error
}

func (r *ExerciseRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Exercise, error) {
	var exercise domain.Exercise
	err := r.DB.WithContext(ctx).First(&exercise, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}

func (r *ExerciseRepo) GetByLesson(ctx context.Context, lessonID uuid.UUID) (*[]domain.Exercise, error) {
	var exercises []domain.Exercise
	err := r.DB.WithContext(ctx).Where("lesson_id = ?", lessonID).Find(&exercises).Error
	return &exercises, err
}

func (r *ExerciseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&domain.Exercise{}, "id = ?", id).Error
}
