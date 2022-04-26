package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS" envDefault:":8080"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseURL          string `env:"DATABASE_URI" envDefault:"postgresql://localhost:5432/product?user=postgres&password=password"`
	SecretKey            string
	MigrationDir         string
}

const (
	defaultServerAddress = ":8080"
	defaultAccrualAdd    = ""
	defaultDB            = "postgresql://localhost:5432/product?user=postgres&password=password"
	defaultSecretKey     = "SecretKey"
	defaultmigration     = "./migrations"
)

var defaultConfig = Config{
	RunAddress:           defaultServerAddress,
	AccrualSystemAddress: defaultAccrualAdd,
	DatabaseURL:          defaultDB,
	SecretKey:            defaultSecretKey,
}

func NewConfig() (Config, error) {
	conf := defaultConfig
	conf.parseFlags()
	conf.parseEnvVars()
	err := conf.Validate()
	return conf, err
}

func (conf *Config) parseFlags() {

	flag.StringVar(&conf.RunAddress, "a", defaultServerAddress, "network address the server listens on")
	flag.StringVar(&conf.AccrualSystemAddress, "r", defaultAccrualAdd, "Accrual system address")
	flag.StringVar(&conf.DatabaseURL, "d", defaultDB, "database")

	flag.Parse()

}

func (conf *Config) parseEnvVars() {

	ra := os.Getenv("RUN_ADDRESS")
	if ra != "" {
		conf.RunAddress = ra
	}

	asa, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if ok {

		conf.AccrualSystemAddress = asa
	}

	db, ok := os.LookupEnv("DATABASE_URI")
	if ok {

		conf.DatabaseURL = db
	}
}

func (conf *Config) Validate() error {

	conf.RunAddress = strings.TrimSpace(conf.RunAddress)

	return nil
}
