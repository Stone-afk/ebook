package ratelimit

import (
	"context"
	"errors"
	"sync"
	"time"
)

// LeakyBucket 漏桶算法
type LeakyBucket struct {
	// 每隔多久一个令牌
	interval time.Duration

	closeCh   chan struct{}
	closeOnce sync.Once
	ticker    *time.Ticker
}

func (l *LeakyBucket) Limit(ctx context.Context, key string) (bool, error) {
	select {
	case <-l.ticker.C:
		// 拿到了令牌
		return false, nil
	case <-ctx.Done():
		return true, ctx.Err()
	case <-l.closeCh:
		return true, errors.New("限流器被关了")
	}
}

func (l *LeakyBucket) Close() error {
	l.closeOnce.Do(func() {
		close(l.closeCh)
	})
	return nil
}
