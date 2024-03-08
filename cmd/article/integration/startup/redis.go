package startup

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var redisCmd redis.Cmdable

func InitRedis() redis.Cmdable {
	if redisCmd == nil {
		redisCmd = redis.NewClient(&redis.Options{
			Addr: "localhost:16379",
		})

		for err := redisCmd.Ping(context.Background()).Err(); err != nil; {
			panic(err)
		}
	}
	return redisCmd
}
