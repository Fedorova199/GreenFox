package storage

import (
	"context"
	"database/sql"

	"github.com/Fedorova199/GreenFox/internal/models"
)

type OrderRepository struct {
	db *sql.DB
}

func CreateOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order models.Order) error {
	sqlStatement := `INSERT INTO "order" (number, status, uploaded_at, user_id) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, sqlStatement, order.Number, order.Status, order.UploadedAt, order.UserID)
	return err
}

func (r *OrderRepository) GetByUserID(ctx context.Context, userID uint64) ([]models.Order, error) {
	var orders []models.Order

	rows, err := r.db.QueryContext(ctx, `SELECT id, number, status, accrual, uploaded_at, user_id FROM "order" WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, sql.ErrNoRows
	}

	return orders, nil
}

func (r *OrderRepository) GetByNumber(ctx context.Context, number string) (models.Order, error) {
	var order models.Order

	sqlStatement := `SELECT id, number, status, accrual, uploaded_at, user_id FROM "order" WHERE number = $1`
	row := r.db.QueryRowContext(ctx, sqlStatement, number)
	err := row.Scan(&order.ID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r *OrderRepository) UpdateAccrualStatus(ctx context.Context, accrual models.Accrual) error {
	sqlStatement := `UPDATE "order" SET status = $1, accrual = $2 WHERE number = $3`
	_, err := r.db.ExecContext(ctx, sqlStatement, accrual.Status, accrual.Accrual, accrual.Order)
	return err
}
