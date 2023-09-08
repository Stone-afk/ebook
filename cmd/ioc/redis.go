package ioc

import (
	"ebook/cmd/config"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	redisCmd := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisCmd
}
