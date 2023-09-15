package ratelimit

import (
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed lua/slide_window.lua
var luaSlideWindow string

// RedisSlidingWindowLimiter Redis 上的滑动窗口算法限流器实现
type RedisSlidingWindowLimiter struct {
	cmd redis.Cmdable
	// 窗口大小
	interval time.Duration
	// 阈值
	rate int
	// interval 内允许 rate 个请求
	// 1s 内允许 3000 个请求
}
