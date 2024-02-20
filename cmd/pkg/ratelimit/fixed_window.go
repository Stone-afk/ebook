package ratelimit

import (
	"context"
	"sync"
	"time"
)

// FixedWindowLimiter 固定窗口算法
type FixedWindowLimiter struct {
	// 当前窗口的请求数量
	cnt int
	// 窗口允许的最大的请求数量
	threshold int
	// 窗口大小
	window time.Duration
	// 上一个窗口的起始时间
	lastStart time.Time
	lock      sync.Mutex
}

func (l *FixedWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	l.lock.Lock()
	now := time.Now()
	// 要换窗口了
	if now.After(l.lastStart.Add(l.window)) {
		l.lastStart = now
		l.cnt = 0
	}
	l.cnt++
	if l.cnt > l.threshold {
		l.lock.Unlock()
		return true, nil
	}
	l.lock.Unlock()
	return false, nil
}
