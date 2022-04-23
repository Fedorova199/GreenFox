package storage

import (
	"context"
	"database/sql"

	"github.com/Fedorova199/GreenFox/internal/models"
)

type User struct {
	db *sql.DB
}

func CreateUser(db *sql.DB) *User {
	return &User{
		db: db,
	}
}

func (r *User) Create(ctx context.Context, user models.User) error {
	sqlStatement := `INSERT INTO "user" (login, password_hash) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, sqlStatement, user.Login, user.PasswordHash)
	return err
}

func (r *User) GetByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User

	row := r.db.QueryRowContext(ctx, `SELECT id, login, password_hash, balance, withdrawn FROM "user" WHERE login = $1`, login)
	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.Balance, &user.Withdrawn)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *User) DecreaseBalanceByUserID(ctx context.Context, userID uint64, amount float64) error {
	sqlStatement := `UPDATE "user" SET balance = balance - $1, withdrawn = withdrawn + $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, sqlStatement, amount, userID)
	return err
}
