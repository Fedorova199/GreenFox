package storage

import (
	"context"
	"database/sql"
	"errors"
	"sort"

	"github.com/Fedorova199/GreenFox/internal/models"
	"github.com/Fedorova199/GreenFox/internal/storage/logger"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type Withdrawal interface {
	Create(ctx context.Context, withdrawal models.Withdrawal) error
	GetByUserID(ctx context.Context, userID uint64) ([]models.Withdrawal, error)
}

type WithdrawalDB struct {
	db *sql.DB
}

func CreateWithdrawal(db *sql.DB) *WithdrawalDB {
	return &WithdrawalDB{
		db: db,
	}
}

func (r *WithdrawalDB) Create(ctx context.Context, withdrawal models.Withdrawal) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			logger.Warningf("transaction error: %v", err)
			tx.Rollback()
			return
		}
	}()

	var balance float64
	row := tx.QueryRowContext(ctx, `SELECT balance FROM "user" WHERE id = $1`, withdrawal.UserID)
	err = row.Scan(&balance)
	if err != nil {
		return err
	}
	logger.Debugf("checking the balance: balance %v - required amount %v", balance, withdrawal.Sum)
	if balance < withdrawal.Sum {
		return ErrInsufficientBalance
	}

	createWithdrawalStatement := `INSERT INTO withdrawal ("order", sum, processed_at, user_id) VALUES ($1, $2, $3, $4)`
	_, err = tx.ExecContext(ctx, createWithdrawalStatement, withdrawal.Order, withdrawal.Sum, withdrawal.ProcessedAt, withdrawal.UserID)
	if err != nil {
		return err
	}

	updateBalanceStatement := `UPDATE "user" SET balance = balance - $1, withdrawn = withdrawn + $1 WHERE id = $2`
	_, err = tx.ExecContext(ctx, updateBalanceStatement, withdrawal.Sum, withdrawal.UserID)
	if err != nil {
		return err
	}
	logger.Info(tx.Commit())
	return tx.Commit()
}

func (r *WithdrawalDB) GetByUserID(ctx context.Context, userID uint64) ([]models.Withdrawal, error) {
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

	sort.Slice(withdrawals, func(i, j int) bool {
		return withdrawals[i].ProcessedAt.Before(withdrawals[j].ProcessedAt)
	})

	return withdrawals, nil
}
