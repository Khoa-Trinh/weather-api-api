package cache

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestCache_GetSet(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	c := New(rc, 200*time.Millisecond)

	ctx := context.Background()

	// Miss
	if _, ok, err := c.Get(ctx, "k"); err != nil || ok {
		t.Fatalf("expected miss, got ok=%v err=%v", ok, err)
	}

	// Set then Get
	if err := c.Set(ctx, "k", "v"); err != nil {
		t.Fatal(err)
	}
	if val, ok, err := c.Get(ctx, "k"); err != nil || !ok || val != "v" {
		t.Fatalf("expected hit=v, got ok=%v val=%q err=%v", ok, val, err)
	}

	// Instead of time.Sleep, advance miniredis clock
	mr.FastForward(300 * time.Millisecond)

	if _, ok, err := c.Get(ctx, "k"); err != nil || ok {
		t.Fatalf("expected expired miss, got ok=%v err=%v", ok, err)
	}
}
