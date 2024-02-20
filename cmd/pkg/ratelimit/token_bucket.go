package ratelimit

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

// TokenBucketLimiter 令牌桶算法
type TokenBucketLimiter struct {
	//ch      *time.Ticker
	buckets chan struct{}
	closeCh chan struct{}
	// 每隔多久一个令牌
	interval time.Duration
}

// NewTokenBucketLimiter 把 capacity 设置成0，就是漏桶算法
// 但是，代码可以简化
func NewTokenBucketLimiter(interval time.Duration, capacity int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		interval: interval,
		buckets:  make(chan struct{}, capacity),
	}
}

func (l *TokenBucketLimiter) NewServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		panic("")
	}
}

// Close 你是不能反复调用
func (l *TokenBucketLimiter) Close() error {
	close(l.closeCh)
	return nil
}