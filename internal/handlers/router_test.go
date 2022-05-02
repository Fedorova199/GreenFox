package handlers

import (
	"database/sql"
	"net/http"
	"testing"

	"github.com/Fedorova199/GreenFox/internal/authenticator"
	"github.com/Fedorova199/GreenFox/internal/config"
	middleware "github.com/Fedorova199/GreenFox/internal/middlewares"
	"github.com/Fedorova199/GreenFox/internal/storage"
	"github.com/Fedorova199/GreenFox/internal/storage/logger"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestNewHandler(t *testing.T) {

	cfg := config.ParseVariables()

	db, err := sql.Open("pgx", cfg.DatabasURL)
	if err != nil {
		logger.Fatalln(err)
	}
	defer db.Close()
	userRepository := storage.CreateUser(db)
	orderRepository := storage.CreateOrder(db)
	withdrawalRepository := storage.CreateWithdrawal(db)
	cookieAuthenticator := authenticator.NewCookieAuthenticator([]byte(cfg.SecretKey))
	pointAccrualService := authenticator.NewPointAccrualService(cfg.AccrualSystemAddress, orderRepository)
	pointAccrualService.Start()
	authenticator := middleware.NewAuthenticator(cookieAuthenticator)

	mws := []Middleware{
		middleware.GzipEncoder{},
		middleware.GzipDecoder{},
	}

	handler := NewHandler(
		cfg.AccrualSystemAddress,
		userRepository,
		orderRepository,
		withdrawalRepository,
		cookieAuthenticator,
		pointAccrualService,
		authenticator,
		mws,
	)

	assert.Implements(t, (*http.Handler)(nil), handler)
}
