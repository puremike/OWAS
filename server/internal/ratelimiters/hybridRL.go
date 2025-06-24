package ratelimiters

import "sync"

type HybridLimiter struct {
	SlidingWindow *SlidingWindow
	TokenBucket   *TokenBucket
	mu            sync.Mutex
}

func NewHybridLimiters(slidingWindow *SlidingWindow, tokenBucket *TokenBucket) *HybridLimiter {
	return &HybridLimiter{
		SlidingWindow: slidingWindow,
		TokenBucket:   tokenBucket,
	}
}

func (h *HybridLimiter) Allowed() bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	// check sliding window
	if !h.SlidingWindow.Allowed() {
		return false
	}

	// check token bucket
	return h.TokenBucket.Allowed()
}
