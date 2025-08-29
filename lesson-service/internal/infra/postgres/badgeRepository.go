package postgres

import (
	"lesson-service/internal/domain"

	"gorm.io/gorm"
)

type BadgeRepo struct {
	DB *gorm.DB
}

func NewBadgeRepo(db *gorm.DB) *BadgeRepo {
	return &BadgeRepo{DB: db}
}

func (r *BadgeRepo) GetBadgesByUserID(userID string) (*[]domain.UserBadge, error) {
	var badges []domain.UserBadge
	result := r.DB.Where("user_id = ?", userID).Find(&badges)

	return &badges, result.Error
}

func (r *BadgeRepo) CheckIfBadgeAwarded(userID string, badgeName string)  error {
	var badge domain.UserBadge
	// badge_name, err := badgeName.String()
	// if err != nil {
	// 	return err
	// }

	result := r.DB.First(&badge, "user_id = ? AND badge_name = ?", userID, badgeName)
	// if result.Error != nil {
	// 	if result.Error == gorm.ErrRecordNotFound {
	// 		return nil
	// 	}
	// 	return result.Error
	// }
	// return nil
	return result.Error
}


func (r *BadgeRepo) AddBadge(badge *domain.UserBadge) error {
	return r.DB.Create(badge).Error
}