package config

import (
	"flag"
	"os"
	"strings"
)

type Config struct {
	RUN_ADDRESS            string `env:"RUN_ADDRESS" envDefault:":8080"`
	ACCRUAL_SYSTEM_ADDRESS string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DATABASE_URI           string `env:"DATABASE_URI" envDefault:"postgresql://localhost:5432/test_db?user=postgres&password=password"`
	SecretKey              string
}

const (
	defaultServerAddress = ":8080"
	defaultAccrualAdd    = ""
	defaultDB            = "postgresql://localhost:5432/test_db?user=postgres&password=password"
	defaultSecretKey     = "Secret_Key"
)

var defaultConfig = Config{
	RUN_ADDRESS:            defaultServerAddress,
	ACCRUAL_SYSTEM_ADDRESS: defaultAccrualAdd,
	DATABASE_URI:           defaultDB,
	SecretKey:              defaultSecretKey,
}

func NewConfig() (Config, error) {
	conf := defaultConfig
	conf.parseFlags()
	conf.parseEnvVars()
	err := conf.Validate()
	return conf, err
}

func (conf *Config) parseFlags() {

	flag.StringVar(&conf.RUN_ADDRESS, "a", defaultServerAddress, "network address the server listens on")
	flag.StringVar(&conf.ACCRUAL_SYSTEM_ADDRESS, "r", defaultAccrualAdd, "Accrual system address")
	flag.StringVar(&conf.DATABASE_URI, "d", defaultDB, "database")

	flag.Parse()

}

func (conf *Config) parseEnvVars() {
	ra := os.Getenv("RUN_ADDRESS")
	if ra != "" {
		conf.RUN_ADDRESS = ra
	}

	asa, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if ok {

		conf.ACCRUAL_SYSTEM_ADDRESS = asa
	}

	db, ok := os.LookupEnv("DATABASE_URI")
	if ok {

		conf.DATABASE_URI = db
	}
}

func (conf *Config) Validate() error {

	conf.RUN_ADDRESS = strings.TrimSpace(conf.RUN_ADDRESS)

	return nil
}
