package ratelimit

import (
	"context"
	"github.com/ecodeclub/ekit/queue"
	"google.golang.org/grpc"
	"sync"
	"time"
)

// SlideWindowLimiter 滑动窗口算法
type SlideWindowLimiter struct {
	threshold int
	lock      sync.Mutex
	window    time.Duration
	// 请求的时间戳
	queue queue.PriorityQueue[time.Time]
}

func (l *SlideWindowLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		l.lock.Lock()
		now := time.Now()
		if l.queue.Len() < l.threshold {
			_ = l.queue.Enqueue(time.Now())
			l.lock.Unlock()
			return false, nil
		}
		windowStart := now.Add(-l.window)
		for {
			// 最早的请求
			first, _ := l.queue.Peek()
			if first.After(windowStart) {
				// 退出循环
				break
			}
			// 就是删了 first
			_, _ = l.queue.Dequeue()
		}
		if l.queue.Len() < l.threshold {
			_ = l.queue.Enqueue(time.Now())
			l.lock.Unlock()
			return false, nil
		}
		return true, nil
	}
}
