package postgres

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"lesson-service/internal/domain"
)

type UnitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) *UnitRepository {
	return &UnitRepository{db: db}
}

func (r *UnitRepository) Create(ctx context.Context, unit *domain.Unit) error {
	return r.db.WithContext(ctx).Create(unit).Error
}

func (r *UnitRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Unit, error) {
	var unit domain.Unit
	err := r.db.WithContext(ctx).First(&unit, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *UnitRepository) ListAll(ctx context.Context) (*[]domain.Unit, error) {
	var units []domain.Unit
	err := r.db.WithContext(ctx).Order("order_index ASC").Find(&units).Error
	return &units, err
}

func (r *UnitRepository) Update(ctx context.Context, unit *domain.Unit) error {
	return r.db.WithContext(ctx).Save(unit).Error
}

func (r *UnitRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Unit{}, "id = ?", id).Error
}
