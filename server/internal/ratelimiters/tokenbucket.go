package ratelimiters

import (
	"sync"
	"time"
)

type TokenBucket struct {
	rate           float64    // tokens per second
	capacity       float64    // maximum number of tokens
	mu             sync.Mutex // mutex for thread safety
	tokens         float64    // current number of tokens
	lastRefillTime time.Time  // last time tokens were refilled
}

func NewTokenBucket(rate, capacity float64) *TokenBucket {
	return &TokenBucket{
		rate:           rate,
		capacity:       capacity,
		tokens:         capacity,
		lastRefillTime: time.Now(),
	}
}

func (t *TokenBucket) Allowed() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	elapsedTime := now.Sub(t.lastRefillTime).Seconds()
	t.lastRefillTime = now

	// Refill tokens
	t.tokens += elapsedTime * t.rate
	if t.tokens > t.capacity {
		t.tokens = t.capacity
	}

	// check if token is available
	if t.tokens >= 1 {
		t.tokens--
		return true
	}

	return false
}
