package ioc

import (
	"ebook/cmd/internal/handler"
	"ebook/cmd/internal/handler/middleware"
	"ebook/cmd/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitMiddlewares(redisCmd redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		sessions.Sessions("SESS", memstore.NewStore(handler.SessAuthKey, handler.SessEncryptionKey)),
		middleware.NewJWTLoginMiddlewareBuilder().Build(),
		ratelimit.NewBuilder(redisCmd, time.Second, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	// 跨域拦截器
	return cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 不加这个，前端是拿不到 jwt
		ExposeHeaders: []string{"x-jwt-token"},
		// 是否允许你带 cookie 之类的东西
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "company.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
