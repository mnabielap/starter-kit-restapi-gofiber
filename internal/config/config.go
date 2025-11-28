package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv                      string
	Port                        string
	DBDriver                    string
	DBName                      string
	DBHost                      string
	DBPort                      string
	DBUser                      string
	DBPassword                  string
	DBSSLMode                   string
	JWTSecret                   string
	JWTAccessExpirationMinutes  int
	JWTRefreshExpirationDays    int
	JWTResetPwdExpirationMin    int
	JWTVerifyEmailExpirationMin int
	SMTPHost                    string
	SMTPPort                    int
	SMTPUsername                string
	SMTPPassword                string
	EmailFrom                   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	return &Config{
		AppEnv:                      getEnv("APP_ENV", "development"),
		Port:                        getEnv("PORT", "3000"),
		DBDriver:                    getEnv("DB_DRIVER", "sqlite"),
		DBName:                      getEnv("DB_NAME", "gofiber_app.db"),
		DBHost:                      getEnv("DB_HOST", "localhost"),
		DBPort:                      getEnv("DB_PORT", "5432"),
		DBUser:                      getEnv("DB_USER", "postgres"),
		DBPassword:                  getEnv("DB_PASSWORD", "password"),
		DBSSLMode:                   getEnv("DB_SSLMODE", "disable"),
		JWTSecret:                   getEnv("JWT_SECRET", "secret"),
		JWTAccessExpirationMinutes:  getEnvAsInt("JWT_ACCESS_EXPIRATION_MINUTES", 30),
		JWTRefreshExpirationDays:    getEnvAsInt("JWT_REFRESH_EXPIRATION_DAYS", 30),
		JWTResetPwdExpirationMin:    getEnvAsInt("JWT_RESET_PASSWORD_EXPIRATION_MINUTES", 10),
		JWTVerifyEmailExpirationMin: getEnvAsInt("JWT_VERIFY_EMAIL_EXPIRATION_MINUTES", 10),
		SMTPHost:                    getEnv("SMTP_HOST", ""),
		SMTPPort:                    getEnvAsInt("SMTP_PORT", 587),
		SMTPUsername:                getEnv("SMTP_USERNAME", ""),
		SMTPPassword:                getEnv("SMTP_PASSWORD", ""),
		EmailFrom:                   getEnv("EMAIL_FROM", "support@yourapp.com"),
	}
}

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return def
}

func getEnvAsInt(key string, def int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return def
}