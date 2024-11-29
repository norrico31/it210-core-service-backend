package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost   string
	Port         string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecret    string
	DATABASE_URL string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		PublicHost:   getEnv("", "127.0.0.1"),
		Port:         getEnv("", "8080"),
		DBUser:       getEnv("", "postgres"),
		DBPassword:   getEnv("", "root"),
		DBName:       getEnv("", "it210"),
		JWTSecret:    getEnv("JWT_SECRET", "IS-IT_REALL-A_SECRET-?~JWT-NOT_SO-SURE"),
		DATABASE_URL: getEnv("", "127.0.0.1"),

		// PublicHost:             getEnv("DATABASE_URL", ""),
		// Port:                   getEnv("PORT", "8080"),
		// DBUser:                 getEnv("POSTGRES_USER", ""),
		// DBPassword:             getEnv("PGPASSWORD", ""),
		// DBName:                 getEnv("POSTGRES_DB", ""),
		// JWTSecret:              getEnv("JWT_SECRET", "IS-IT_REALL-A_SECRET-?~JWT-NOT_SO-SURE"),
		// DATABASE_URL:           getEnv("DATABASE_PUBLIC_URL", ""),
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
