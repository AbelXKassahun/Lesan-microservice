package postgres

import (
	"lesson-service/internal/domain"

	"gorm.io/gorm"
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

	return &xp, result.Error
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
	xp.League = domain.DetermineLeague(xp.Total)

	return r.DB.Save(&xp).Error
}


func (r *XPRepo)GetUsersByLeague(league string) (*[]domain.UserXP, error){
	var users []domain.UserXP
    result := r.DB.Where("league = ?", league).Find(&users)
    return &users, result.Error
}
