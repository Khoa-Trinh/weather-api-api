package http

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/Khoa-Trinh/weather-api-api/internal/cache"
	"github.com/Khoa-Trinh/weather-api-api/internal/weather"
)

func silentLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestHealth_OK(t *testing.T) {
	h := &Handler{Log: silentLogger()}
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	h.Health(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if got := rr.Body.String(); got == "" {
		t.Fatalf("expected body, got empty")
	}
}

func TestCacheKey_Normalizes(t *testing.T) {
	k := cacheKey("  Ho Chi MINH  ", " Metric ")
	if k != "wx:ho chi minh:metric" {
		t.Fatalf("got %q", k)
	}
}

func TestGetWeather_RequiresCity(t *testing.T) {
	h := &Handler{Log: silentLogger()}
	req := httptest.NewRequest(http.MethodGet, "/v1/weather", nil)
	rr := httptest.NewRecorder()
	h.GetWeather(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rr.Code)
	}
}

func TestGetWeather_CacheHit_ShortCircuits(t *testing.T) {
	mr, _ := miniredis.Run()
	defer mr.Close()

	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	c := cache.New(rc, 0) // ttl not relevant here

	log := silentLogger()
	// Preload cache with expected key
	key := cacheKey("hanoi", "metric")
	if err := c.Set(context.Background(), key, `{"ok":true}`); err != nil {
		t.Fatalf("preload cache: %v", err)
	}

	h := &Handler{
		Cache:   c,
		Weather: &weather.Client{}, // shouldn't be called on HIT
		Log:     log,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/weather?city=Hanoi&units=metric", nil)
	rr := httptest.NewRecorder()
	h.GetWeather(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rr.Code)
	}
	if got := rr.Header().Get("X-Cache"); got != "HIT" {
		t.Fatalf("want X-Cache=HIT, got %q", got)
	}
	if rr.Body.String() != `{"ok":true}` {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}
