package ioc

import (
	"context"
	"ebook/cmd/internal/handler"
	ijwt "ebook/cmd/internal/handler/jwt"
	"ebook/cmd/internal/handler/middleware"
	loggerMiddleware "ebook/cmd/pkg/ginx/middleware/logger"
	"ebook/cmd/pkg/ginx/middleware/metric"
	limitMiddleware "ebook/cmd/pkg/ginx/middleware/ratelimit"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/ratelimit"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *handler.UserHandler,
	oauth2WechatHdl *handler.OAuth2WechatHandler, articleHdl *handler.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	oauth2WechatHdl.RegisterRoutes(server)
	(&handler.ObservabilityHandler{}).RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisCmd redis.Cmdable, jwtHdl ijwt.Handler, l logger.Logger) []gin.HandlerFunc {
	limiter := ratelimit.NewRedisSlidingWindowLimiter(redisCmd, time.Second, 100)
	bd := loggerMiddleware.NewBuilder(func(ctx context.Context, al *loggerMiddleware.AccessLog) {
		l.Debug("HTTP请求", logger.Field{Key: "AccessLog", Value: al})
	})
	viper.OnConfigChange(func(in fsnotify.Event) {
		ok := viper.GetBool("web.logReqBody")
		bd.AllowReqBody(ok)
		ok = viper.GetBool("web.logRespBody")
		bd.AllowRespBody(ok)
	})
	mb := &metric.MiddlewareBuilder{
		Namespace:  "web server",
		Subsystem:  "ebook",
		Name:       "gin_http",
		Help:       "统计 GIN 的 HTTP 接口",
		InstanceID: "my-instance-1",
	}
	return []gin.HandlerFunc{
		corsHdl(),
		mb.Build(),
		bd.Build(),
		sessions.Sessions("SESS", memstore.NewStore(handler.SessAuthKey, handler.SessEncryptionKey)),
		middleware.NewJWTLoginMiddlewareBuilder(jwtHdl).Build(),
		limitMiddleware.NewBuilder(limiter).Build(),
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
