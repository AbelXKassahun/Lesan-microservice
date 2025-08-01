package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"gamification-service/internal/models"
)

type XPRepo struct {
	DB *pgx.Conn
}

func NewXPRepo(db *pgx.Conn) *XPRepo {
	return &XPRepo{DB: db}
}

func (r *XPRepo) GetXPByUserID(userID string) (*models.UserXP, error) {
	var xp models.UserXP
	err := r.DB.QueryRow(context.Background(),
		`SELECT user_id, total_xp FROM user_xp WHERE user_id = $1`, userID,
	).Scan(&xp.UserID, &xp.Total)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &xp, err
}

func (r *XPRepo) AddXP(userID string, amount int) error {
	_, err := r.DB.Exec(context.Background(), `
		INSERT INTO user_xp (user_id, total_xp)
		VALUES ($1, $2)
		ON CONFLICT (user_id)
		DO UPDATE SET total_xp = user_xp.total_xp + $2
	`, userID, amount)

	return err
}
