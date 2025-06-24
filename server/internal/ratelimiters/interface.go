package ratelimiters

type Limiter interface {
	Allowed() bool
}

var _ Limiter = (*HybridLimiter)(nil)
