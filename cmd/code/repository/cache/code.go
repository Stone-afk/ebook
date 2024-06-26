package cache

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

// 编译器会在编译的时候，把 set_code 的代码放进来这个 luaSetCode 变量里
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

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

func (rc *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := rc.cmd.Eval(ctx, luaSetCode, []string{rc.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		// 毫无问题
		return nil
	case 1:
		// 发送太频繁
		return ErrCodeSendTooMany
	//case -2:
	//	return
	default:
		// 系统错误
		return errors.New("系统错误")
	}
}

func (rc *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := rc.cmd.Eval(ctx, luaVerifyCode, []string{rc.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		// 正常来说，如果频繁出现这个错误，你就要告警，因为有人搞你
		return false, ErrCodeVerifyTooManyTimes
	case -2:
		return false, nil
		//default:
		//	return false, ErrUnknownForCode
	}
	return false, ErrUnknownForCode
}

func (rc *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
