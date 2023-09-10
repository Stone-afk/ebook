package main

import (
	"ebook/cmd/config"
	"ebook/cmd/internal/handler"
	"ebook/cmd/internal/handler/middleware"
	"ebook/cmd/internal/repository"
	"ebook/cmd/internal/repository/cache"
	"ebook/cmd/internal/repository/dao"
	"ebook/cmd/internal/service"
	"ebook/cmd/internal/service/sms/memory"
	"ebook/cmd/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {
	db := initDB()
	rc := initRedis()
	u := initUser(db, rc)

	server := initServer()
	u.RegisterRoutes(server)

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})

	err := server.Run(":8083")
	if err != nil {
		panic(err)
	}
}

func initServer() *gin.Engine {
	server := gin.Default()

	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	// 跨域拦截器
	server.Use(cors.New(cors.Config{
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
	}))

	// 设置 session
	//// 设置到 cookie
	//store := cookie.NewStore([]byte("secret"))
	//// 设置到本地内存
	store := memstore.NewStore(handler.SessAuthKey, handler.SessEncryptionKey)
	// 设置到redis
	//store, err := redisSess.NewStore(16,
	//	"tcp", "localhost:6379", "",
	//	handler.SessAuthKey, handler.SessEncryptionKey)
	//
	//if err != nil {
	//	panic(err)
	//}
	server.Use(sessions.Sessions("SESS", store))

	// jwt 登录中间件
	server.Use(middleware.NewJWTLoginMiddlewareBuilder().Build())

	return server

}

func initUser(db *gorm.DB, rcCmd redis.Cmdable) *handler.UserHandler {
	rc := cache.NewRedisUserCache(rcCmd)
	codeRc := cache.NewCodeCache(rcCmd)
	ud := dao.NewGORMUserDAO(db)
	repo := repository.NewUserRepository(ud, rc)
	codeRepo := repository.NewCodeRepository(codeRc)
	svc := service.NewUserService(repo)
	codeSvc := service.NewCodeService(codeRepo, memory.NewService())
	u := handler.NewUserHandler(svc, codeSvc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
