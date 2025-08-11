package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port             string
	VCAPIKey         string
	DefaultUnitGroup string
	RedisAddr        string
	RedisPassword    string
	RedisDB          int
	CacheTTL         time.Duration
	HTTPTimeout      time.Duration
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env %s", key)
	}
	return v
}

func parseInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			return i
		}
	}
	return def
}

func Load() Config {
	return Config{
		Port:             getEnv("PORT", "8080"),
		VCAPIKey:         mustEnv("VISUAL_CROSSING_API_KEY"),
		DefaultUnitGroup: getEnv("DEFAULT_UNIT_GROUP", "metric"),
		RedisAddr:        getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          parseInt("REDIS_DB", 0),
		CacheTTL:         time.Duration(parseInt("CACHE_TTL_SECONDS", 43200)) * time.Second,
		HTTPTimeout:      time.Duration(parseInt("HTTP_TIMEOUT_SECONDS", 10)) * time.Second,
	}
}
