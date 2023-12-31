package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	addr := viper.GetString("redis.addr")
	redisCmd := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return redisCmd
}

func InitRLockClient(cmd redis.Cmdable) *rlock.Client {
	return rlock.NewClient(cmd)
}
