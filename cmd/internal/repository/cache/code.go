package cache

import (
	"context"
	_ "embed"
	"errors"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
	ErrUnknownForCode         = errors.New("我也不知发生什么了，反正是跟 code 有关")
)

// 编译器会在编译的时候，把 set_code 的代码放进来这个 luaSetCode 变量里
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	cmd redis.Cmdable
}

// NewCodeCacheGoBestPractice Go 的最佳实践是返回具体类型
func NewCodeCacheGoBestPractice(cmd redis.Cmdable) *RedisCodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

func NewCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	panic("")
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	panic("")
}
