package main

import (
	"database/sql"
	"fmt"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Fedorova199/GreenFox/internal/authenticator"
	"github.com/Fedorova199/GreenFox/internal/config"
	"github.com/Fedorova199/GreenFox/internal/handlers"
	middleware "github.com/Fedorova199/GreenFox/internal/middlewares"
	"github.com/Fedorova199/GreenFox/internal/storage"
	"github.com/Fedorova199/GreenFox/internal/storage/logger"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"

	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	logger.SetLevel("WARNING")
	cfg := config.ParseVariables()
	logger.Debugf("cfg: %v", cfg)
	db, err := sql.Open("pgx", cfg.DatabasURL)
	if err != nil {
		logger.Fatalln(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("could not start sql migration... %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", cfg.MigrationDir), "product", driver)
	if err != nil {
		log.Fatalf("migration failed... %v", err)
	}

	fmt.Println(m.Version())
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}

	userRepository := storage.CreateUser(db)
	orderRepository := storage.CreateOrder(db)
	withdrawalRepository := storage.CreateWithdrawal(db)
	cookieAuthenticator := authenticator.NewCookieAuthenticator([]byte(cfg.SecretKey))
	pointAccrualService := authenticator.NewPointAccrualService(cfg.AccrualSystemAddress, orderRepository)
	pointAccrualService.Start()
	authenticator := middleware.NewAuthenticator(cookieAuthenticator)

	mws := []handlers.Middleware{
		middleware.GzipEncoder{},
		middleware.GzipDecoder{},
	}

	handler := handlers.NewHandler(
		cfg.AccrualSystemAddress,
		userRepository,
		orderRepository,
		withdrawalRepository,
		cookieAuthenticator,
		pointAccrualService,
		authenticator,
		mws,
	)
	server := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: handler,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-c
		server.Close()
		pointAccrualService.Stop()
	}()

	log.Fatal(server.ListenAndServe())
}
