package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"lesson-service/internal/domain"
)

type LessonRepository struct {
	DB *gorm.DB
}

func NewLessonRepository(db *gorm.DB) *LessonRepository {
	return &LessonRepository{DB: db}
}

func (r *LessonRepository) Create(ctx context.Context, lesson *domain.Lesson) error {
	return r.DB.WithContext(ctx).Create(lesson).Error
}

func (r *LessonRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Lesson, error) {
	var lesson domain.Lesson
	err := r.DB.WithContext(ctx).First(&lesson, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (r *LessonRepository) GetByUnit(ctx context.Context, unitID uuid.UUID) (*[]domain.Lesson, error) {
	var lessons []domain.Lesson
	err := r.DB.WithContext(ctx).Where("unit_id = ?", unitID).Order("order_index ASC").Find(&lessons).Error
	return &lessons, err
}

func (r *LessonRepository) GetBySection(ctx context.Context, sectionID uuid.UUID) (*[]domain.Lesson, error) {
	var lessons []domain.Lesson
	err := r.DB.WithContext(ctx).Where("section_id = ?", sectionID).Order("order_index ASC").Find(&lessons).Error
	return &lessons, err
}

func (r *LessonRepository) ListAll(ctx context.Context) (*[]domain.Lesson, error) {
	var lessons []domain.Lesson
	err := r.DB.WithContext(ctx).Order("created_at DESC").Find(&lessons).Error
	return &lessons, err
}

func (r *LessonRepository) Update(ctx context.Context, lesson *domain.Lesson) error {
	return r.DB.WithContext(ctx).Save(lesson).Error
}

func (r *LessonRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&domain.Lesson{}, "id = ?", id).Error
}
