package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
	GatewayPort            string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:             getEnv("DATABASE_URL", "http://localhost"),
		Port:                   getEnv("PORT", "8081"),
		DBUser:                 getEnv("POSTGRES_USER", "postgres"),
		DBPassword:             getEnv("POSTGRES_PASSWORD", ""),
		DBAddress:              getEnv("DB_ADDRESS", "postgres"),
		GatewayPort:            getEnv("GATEWAY_SERVICE_PORT", "8080"),
		DBName:                 getEnv("POSTGRES_DB", "railway"),
		JWTSecret:              getEnv("JWT_SECRET", "IS-IT_REALL-A_SECRET-?~JWT-NOT_SO-SURE"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXP", 3600*24*7),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		envVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return envVal
	}
	return fallback
}
