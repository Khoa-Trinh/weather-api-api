package weather

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Fetch_Success(t *testing.T) {
	// Fake Visual Crossing server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Expect /timeline/<place> with some query params
		if r.URL.Path == "/timeline/Hanoi" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"resolvedAddress":"Hanoi, Vietnam","days":[{"temp":30.5}]}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer ts.Close()

	c := NewClient("dummy-key", "metric", 5*time.Second)
	// Point baseURL to our test server
	c.baseURL = ts.URL + "/timeline"
	// Use test server client (no timeout issues)
	c.httpClient = ts.Client()

	body, status, err := c.Fetch(context.Background(), "Hanoi", "metric")
	if err != nil || status != http.StatusOK {
		t.Fatalf("unexpected err=%v status=%d", err, status)
	}

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("bad json: %v", err)
	}
	if got["resolvedAddress"] != "Hanoi, Vietnam" {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestClient_Fetch_UpstreamError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"message":"bad place"}`, http.StatusBadRequest)
	}))
	defer ts.Close()

	c := NewClient("dummy-key", "metric", 5*time.Second)
	c.baseURL = ts.URL + "/timeline"
	c.httpClient = ts.Client()

	_, status, err := c.Fetch(context.Background(), "NoWhere", "metric")
	if err == nil || status != http.StatusBadRequest {
		t.Fatalf("expected upstream error, got err=%v status=%d", err, status)
	}
}
