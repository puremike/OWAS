package ratelimiters

import (
	"container/ring"
	"sync"
	"time"
)

type SlidingWindow struct {
	limit      int
	window     time.Duration
	mu         sync.RWMutex
	requestLog *ring.Ring
}

func NewSlidingWindow(limit int, window time.Duration) *SlidingWindow {
	return &SlidingWindow{
		limit:      limit,
		window:     window,
		requestLog: ring.New(limit),
	}
}

func (s *SlidingWindow) Allowed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	cutoff := now.Add(-s.window)
	count := 0

	s.requestLog.Do(func(p any) {
		if p != nil && p.(time.Time).After(cutoff) {
			count++
		}
	})

	if count >= s.limit {
		return false // limit exceeded
	}

	s.requestLog.Value = now
	s.requestLog = s.requestLog.Next()
	return true
}
