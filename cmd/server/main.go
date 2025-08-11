package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"

	"github.com/Khoa-Trinh/weather-api-api/internal/cache"
	"github.com/Khoa-Trinh/weather-api-api/internal/config"
	httpx "github.com/Khoa-Trinh/weather-api-api/internal/http"
	"github.com/Khoa-Trinh/weather-api-api/internal/weather"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rc.Ping(ctx).Err(); err != nil {
		logger.Warn("redis unreachable, continuing without cache", slog.String("err", err.Error()))
	}
	var c *cache.Cache
	if err := rc.Ping(ctx).Err(); err == nil {
		c = cache.New(rc, cfg.CacheTTL)
	}

	wc := weather.NewClient(cfg.VCAPIKey, cfg.DefaultUnitGroup, cfg.HTTPTimeout)
	h := httpx.NewHandler(c, wc, logger)

	r := chi.NewRouter()
	for _, mw := range httpx.CommonMiddleware() {
		r.Use(mw)
	}

	r.Get("/health", h.Health)
	r.Get("/weather", h.GetWeather)
	r.Get("/_cache", h.RedisPing)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("starting server", slog.String("port", cfg.Port))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
