package config

import (
	"os"
)

type Config struct {
	MongoURI  string
	DBName    string
	Port      string
	JWTSecret string
}

func LoadConfig() *Config {
	return &Config{
		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DBName:    getEnv("DB_NAME", "echelon_db"),
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
