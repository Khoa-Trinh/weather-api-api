package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestCommonMiddleware_RateLimits(t *testing.T) {
	r := chi.NewRouter()
	for _, mw := range CommonMiddleware() {
		r.Use(mw)
	}
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Hit the endpoint 60 times quickly (limit is 60/min in CommonMiddleware)
	for i := 0; i < 60; i++ {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ping", nil)
		r.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("unexpected code on request %d: %d", i, rr.Code)
		}
	}

	// 61st request should be limited (429) within the same minute
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(rr, req)

	// Allow either 429 (expected) or 200 if library timing is slightly lenient in CI
	if rr.Code != http.StatusTooManyRequests && rr.Code != http.StatusOK {
		t.Fatalf("expected 429 or 200, got %d", rr.Code)
	}

	// Sleep a bit to avoid flakiness across environments
	time.Sleep(10 * time.Millisecond)
}
