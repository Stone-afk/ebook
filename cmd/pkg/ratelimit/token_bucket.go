package ratelimit

import (
	"context"
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
	res := &TokenBucketLimiter{
		interval: interval,
		buckets:  make(chan struct{}, capacity),
	}
	ticker := time.NewTicker(res.interval)
	go func() {
		for {
			select {
			case <-res.closeCh:
				return
			case <-ticker.C:
				select {
				case res.buckets <- struct{}{}:
				default:
				}
			}
		}
		//for _ = range ticker.C {
		//	select {
		//	case l.buckets <- struct{}{}:
		//	// 发到了桶里面
		//	default:
		//		// 桶满了
		//	}
		//}
	}()
	return res
}

func (l *TokenBucketLimiter) Limit(ctx context.Context, key string) (bool, error) {
	select {
	case <-l.buckets:
		// 拿到了令牌
		return false, nil
	case <-ctx.Done():
		// 没有令牌就等令牌，直到超时
		return false, ctx.Err()
	default:
		// 就意味着你认为，没有令牌不应阻塞，直接返回
		return false, nil
	}
}

// Close 你是不能反复调用
func (l *TokenBucketLimiter) Close() error {
	close(l.closeCh)
	return nil
}
