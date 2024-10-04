package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	ConnectionString     string
	JWTExpirationSeconds string
	JWTSecret            string
}

var Envs = initConfig()

func initConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("missing or empty .env: %s\n", err.Error())
	}

	var (
		port                 = os.Getenv("PORT")
		JWTExpirationSeconds = os.Getenv("JWT_EXP")
		connStr              = os.Getenv("DATABASE_URL")
		JWTSecret            = os.Getenv("JWTSecret")
	)

	return Config{
		Port:                 port,
		ConnectionString:     connStr,
		JWTExpirationSeconds: JWTExpirationSeconds,
		JWTSecret:            JWTSecret,
	}
}
