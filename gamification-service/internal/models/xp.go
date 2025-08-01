package models

type UserXP struct {
	UserID string `db:"user_id"`
	Total  int    `db:"total_xp"`
}
