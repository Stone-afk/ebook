package ratelimit

import (
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
