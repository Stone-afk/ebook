package ratelimit

import (
	"context"
	"sync"
	"time"
)

// LeakyBucket 漏桶算法
type LeakyBucket struct {
	// 每隔多久一个令牌
	interval time.Duration

	closeCh   chan struct{}
	closeOnce sync.Once
}

func (l *LeakyBucket) Limit(ctx context.Context, key string) (bool, error) {
	panic("")
}

func (l *LeakyBucket) Close() error {
	l.closeOnce.Do(func() {
		close(l.closeCh)
	})
	return nil
}
