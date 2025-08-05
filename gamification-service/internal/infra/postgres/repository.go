package postgres

import (
	
	"gorm.io/gorm"
	"gamification-service/internal/domain"
)

type XPRepo struct {
	DB *gorm.DB
}

func NewXPRepo(db *gorm.DB) *XPRepo {
	return &XPRepo{DB: db}
}

func (r *XPRepo) GetXPByUserID(userID string) (*domain.UserXP, error) {
	var xp domain.UserXP
	result := r.DB.First(&xp, "user_id = ?", userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &xp, nil
}

func (r *XPRepo) AddXP(userID string, amount int) error {
	var xp domain.UserXP
	result := r.DB.First(&xp, "user_id = ?", userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			xp = domain.UserXP{UserID: userID, Total: amount}
			return r.DB.Create(&xp).Error
		}
		return result.Error
	}

	xp.Total += amount
	return r.DB.Save(&xp).Error
}
