package ratelimit

import (
	"context"
	"sync/atomic"
)

type CounterLimiter struct {
	cnt       *atomic.Int32
	threshold int32
}

func NewCounterLimiter(cnt *atomic.Int32, threshold int32) Limiter {
	return &CounterLimiter{
		cnt:       cnt,
		threshold: threshold,
	}
}

func (l *CounterLimiter) Limit(ctx context.Context, key string) (bool, error) {
	cnt := l.cnt.Add(1)
	defer func() {
		l.cnt.Add(-1)
	}()
	if cnt > l.threshold {
		// 这里就是拒绝
		return true, nil
	}
	return false, nil
}
