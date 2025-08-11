package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/Khoa-Trinh/weather-api-api/internal/cache"
	"github.com/Khoa-Trinh/weather-api-api/internal/weather"
)

type Handler struct {
	Cache   *cache.Cache
	Weather *weather.Client
	Log     *slog.Logger
}

func NewHandler(c *cache.Cache, w *weather.Client, log *slog.Logger) *Handler {
	return &Handler{Cache: c, Weather: w, Log: log}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func cacheKey(city, units string) string {
	city = strings.ToLower(strings.TrimSpace(city))
	units = strings.ToLower(strings.TrimSpace(units))
	return fmt.Sprintf("wx:%s:%s", city, units)
}

func (h *Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	city := r.URL.Query().Get("city")
	units := r.URL.Query().Get("units")
	if strings.TrimSpace(city) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"status": "bad request"})
		return
	}

	key := cacheKey(city, units)

	if h.Cache != nil {
		if val, ok, err := h.Cache.Get(ctx, key); err == nil && ok {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "HIT")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(val))
			return
		}
	}

	body, status, err := h.Weather.Fetch(ctx, city, units)
	if err != nil && status == 0 {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": "upstream_unavailable"})
		return
	}

	if status >= 200 && status < 300 && h.Cache != nil {
		if err := h.Cache.Set(context.Background(), key, string(body)); err != nil {
			h.Log.Warn("cache set failed", slog.String("key", key), slog.String("err", err.Error()))
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if status >= 200 && status < 300 {
		w.Header().Set("X-Cache", "MISS")
	}
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func (h *Handler) RedisPing(w http.ResponseWriter, r *http.Request) {
	if h.Cache == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "cache_disabled"})
		return
	}
	type pinger interface {
		Ping(ctx context.Context) *redis.StatusCmd
	}
	writeJSON(w, http.StatusOK, map[string]string{"cache": "enabled"})
}
