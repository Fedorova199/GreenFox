package main

import (
	"database/sql"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Fedorova199/GreenFox/internal/config"
	"github.com/Fedorova199/GreenFox/internal/handlers"
	"github.com/Fedorova199/GreenFox/internal/interfaces"
	middleware "github.com/Fedorova199/GreenFox/internal/middlewares"
	"github.com/Fedorova199/GreenFox/internal/service"
	"github.com/Fedorova199/GreenFox/internal/storage"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping DB... %v", err)
	}

	user := storage.CreateUser(db)
	order := storage.CreateOrderRepository(db)
	withdrawal := storage.CreateWithdrawalRepository(db)
	cookieAuthenticator := service.NewCookieAuthenticator([]byte(cfg.SecretKey))
	pointAccrualService := service.NewPointAccrualService(cfg.AccrualSystemAddress, order, user)
	pointAccrualService.Start()
	authenticator := middleware.NewAuthenticator(cookieAuthenticator)

	mws := []interfaces.Middleware{
		middleware.GzipEncoder{},
		middleware.GzipDecoder{},
	}

	handler := handlers.NewHandler(
		cfg.AccrualSystemAddress,
		user,
		order,
		withdrawal,
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
