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
	ticker := time.NewTicker(l.interval)
	go func() {
		for {
			select {
			case <-l.closeCh:
				return
			case <-ticker.C:
				select {
				case l.buckets <- struct{}{}:
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
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		select {
		case <-l.buckets:
			// 拿到了令牌
			return handler(ctx, req)
		case <-ctx.Done():
			// 没有令牌就等令牌，直到超时
			return nil, ctx.Err()
			// default:
			// 就意味着你认为，没有令牌不应阻塞，直接返回
			//return nil, status.Errorf(codes.ResourceExhausted, "限流了")
		}
	}
}

// Close 你是不能反复调用
func (l *TokenBucketLimiter) Close() error {
	close(l.closeCh)
	return nil
}
