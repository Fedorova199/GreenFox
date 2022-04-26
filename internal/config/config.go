package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabasURL           string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SecretKey            string
	MigrationDir         string
}

func ParseVariables() Config {
	var cfg = Config{
		RunAddress:           "localhost:8080",
		DatabasURL:           "postgresql://localhost:5432/product?user=postgres&password=password",
		AccrualSystemAddress: "",
		SecretKey:            "SecretKey",
		MigrationDir:         "./migrations",
	}

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "Run address")
	flag.StringVar(&cfg.DatabasURL, "d", cfg.DatabasURL, "Database URI")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", cfg.AccrualSystemAddress, "Accrual system address")
	flag.Parse()

	return cfg
}
