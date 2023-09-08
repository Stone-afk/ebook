package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}

func corsHdl() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}