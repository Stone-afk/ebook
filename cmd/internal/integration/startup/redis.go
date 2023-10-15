package startup

import (
	"github.com/redis/go-redis/v9"
)

var redisCmd redis.Cmdable

func InitRedis() redis.Cmdable {
	if redisCmd == nil {
		redisCmd = redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		})
	}
	return redisCmd
}
