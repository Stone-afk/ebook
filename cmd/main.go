package cmd

import (
	"ebook/cmd/internal/handler"
	"ebook/cmd/internal/repository"
	"ebook/cmd/internal/repository/cache"
	"ebook/cmd/internal/repository/dao"
	"ebook/cmd/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	rc := initRedis()
	u := initUser(db, rc)

	server := initServer()
	u.RegisterRoutes(server)

	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func initServer() *gin.Engine {
	server := gin.Default()

	// 跨域拦截器
	server.Use(cors.New(cors.Config{
		//AllowOrigins: []string{"*"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		//ExposeHeaders: []string{"x-jwt-token"},
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

	// 步骤1
	store := cookie.NewStore([]byte("secret"))
	server.Use(sessions.Sessions("mysession", store))

	// jwt 登录中间件
	server.Use()

	return server

}

func initUser(db *gorm.DB, rcCmd redis.Cmdable) *handler.UserHandler {
	rc := cache.NewRedisUserCache(rcCmd)
	ud := dao.NewGORMUserDAO(db)
	repo := repository.NewUserRepository(ud, rc)
	svc := service.NewUserService(repo)
	u := handler.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/ebook"))
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
		Addr:     "localhost:16379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
