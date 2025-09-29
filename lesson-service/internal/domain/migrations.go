package domain

import (
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Exercise{},
		&Lesson{},
		&Section{},
		&Unit{},
		&UserProgress{},
	)
}