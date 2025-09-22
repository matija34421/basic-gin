package config

import (
	"log"
	"net"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string

	PostgresDSN string

	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMode string

	RedisAddr  string
	RedisPass  string
	HMACSecret string
}

var App Config

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func Load() {
	_ = godotenv.Load()

	App = Config{
		ServerPort: getenv("SERVER_PORT", "8080"),

		DBHost:    getenv("DB_HOST", "localhost"),
		DBPort:    getenv("DB_PORT", "5432"),
		DBUser:    getenv("DB_USER", "postgres"),
		DBPass:    getenv("DB_PASS", "postgres"),
		DBName:    getenv("DB_NAME", "basic_gin"),
		DBSSLMode: getenv("DB_SSLMODE", "disable"),

		RedisAddr:  getenv("REDIS_ADDR", "localhost:6379"),
		RedisPass:  getenv("REDIS_PASS", ""),
		HMACSecret: getenv("HMAC_SECRET", "dev-secret"),

		PostgresDSN: getenv("POSTGRES_DSN", ""),
	}

	if strings.TrimSpace(App.PostgresDSN) == "" {
		App.PostgresDSN = buildDSN(App.DBHost, App.DBPort, App.DBUser, App.DBPass, App.DBName, App.DBSSLMode)
	}

	log.Printf("config loaded: port=%s", App.ServerPort)
}

func buildDSN(host, port, user, pass, dbname, sslmode string) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + dbname,
	}
	q := url.Values{}
	if sslmode != "" {
		q.Set("sslmode", sslmode)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
