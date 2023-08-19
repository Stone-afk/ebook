package ratelimit

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

type Builder struct {
	prefix string
	// 阈值
	rate     int
	cmd      redis.Cmdable
	interval time.Duration
}

//go:embed slide_window.lua
var luaScript string

func NewBuilder(cmd redis.Cmdable, interval time.Duration, rate int) *Builder {
	return &Builder{
		prefix:   "ip-limiter",
		rate:     rate,
		cmd:      cmd,
		interval: interval,
	}
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}

func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	panic("")
}
