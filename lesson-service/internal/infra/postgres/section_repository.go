package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"lesson-service/internal/domain"
)


type SectionRepository struct {
	DB *gorm.DB
}

func NewSectionRepository(db *gorm.DB) *SectionRepository {
	return &SectionRepository{DB: db}
}

func (r *SectionRepository) Create(ctx context.Context, service *domain.Section) error {
	return r.DB.WithContext(ctx).Create(service).Error
}

func (r *SectionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Section, error) {
	var section domain.Section
	err := r.DB.WithContext(ctx).First(&section, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

func (r *SectionRepository) GetByUnit(ctx context.Context, unitID uuid.UUID) (*[]domain.Section, error) {
	var sections []domain.Section
	err := r.DB.WithContext(ctx).Where("unit_id = ?", unitID).Order("order_index ASC").Find(&sections).Error
	return &sections, err
}

func (r *SectionRepository) ListAll(ctx context.Context) (*[]domain.Section, error) {
	var sections []domain.Section
	err := r.DB.WithContext(ctx).Order("created_at DESC").Find(&sections).Error
	return &sections, err
}

func (r *SectionRepository) Update(ctx context.Context, section *domain.Section) error {
	return r.DB.WithContext(ctx).Save(section).Error
}

func (r *SectionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&domain.Section{}, "id = ?", id).Error
}
