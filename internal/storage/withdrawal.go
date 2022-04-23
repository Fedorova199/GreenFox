package storage

import (
	"context"
	"database/sql"

	"github.com/Fedorova199/GreenFox/internal/models"
)

type WithdrawalRepository struct {
	db *sql.DB
}

func CreateWithdrawalRepository(db *sql.DB) *WithdrawalRepository {
	return &WithdrawalRepository{
		db: db,
	}
}

func (r *WithdrawalRepository) Create(ctx context.Context, withdrawal models.Withdrawal) error {
	sqlStatement := `INSERT INTO withdrawal ("order", sum, processed_at, user_id) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, sqlStatement, withdrawal.Order, withdrawal.Sum, withdrawal.ProcessedAt, withdrawal.UserID)
	return err
}

func (r *WithdrawalRepository) GetByUserID(ctx context.Context, userID uint64) ([]models.Withdrawal, error) {
	var withdrawals []models.Withdrawal

	rows, err := r.db.QueryContext(ctx, `SELECT id, "order", sum, processed_at, user_id FROM withdrawal WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var withdrawal models.Withdrawal
		err := rows.Scan(&withdrawal.ID, &withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt, &withdrawal.UserID)
		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, sql.ErrNoRows
	}

	return withdrawals, nil
}
