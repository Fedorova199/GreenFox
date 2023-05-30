package storage

import (
	"context"
	"database/sql"

	"github.com/Fedorova199/GreenFox/internal/models"
)

type User interface {
	Create(ctx context.Context, user models.User) error
	GetByLogin(ctx context.Context, login string) (models.User, error)
}
type UserDB struct {
	db *sql.DB
}

func CreateUser(db *sql.DB) *UserDB {
	return &UserDB{
		db: db,
	}
}

func (r *UserDB) Create(ctx context.Context, user models.User) error {
	sqlStatement := `INSERT INTO "user" (login, password_hash) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, sqlStatement, user.Login, user.PasswordHash)
	return err
}

func (r *UserDB) GetByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User
	row := r.db.QueryRowContext(ctx, `SELECT id, login, password_hash, balance, withdrawn FROM "user" WHERE login = $1`, login)
	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.Balance, &user.Withdrawn)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
