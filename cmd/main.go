package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	//db := initDB()
	//rc := initRedis()
	//u := initUser(db, rc)
	//
	//server := initServer()
	//u.RegisterRoutes(server)
	initViperWatch()
	//initLogger()
	keys := viper.AllKeys()
	println(keys)
	setting := viper.AllSettings()
	fmt.Println(setting)

	server := InitWebServer()

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})

	err := server.Run(":8083")
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.L().Info("这是 replace 之前")
	// 如果你不 replace，直接用 zap.L()，你啥都打不出来。
	zap.ReplaceGlobals(logger)
	zap.L().Info("hello，你搞好了")
	type Demo struct {
		Name string `json:"name"`
	}
	zap.L().Info("这是实验参数",
		zap.Error(errors.New("这是一个 error")),
		zap.Int64("id", 123),
		zap.Any("一个结构体", Demo{Name: "hello"}))
}

func initViper() {
	viper.SetDefault("db.dsn", "root:root@tcp(localhost:3306)/mysql")
	// 读取的文件名字叫做 dev
	viper.SetConfigName("dev")
	// 读取的类型是 yaml 文件
	viper.SetConfigType("yaml")
	// 在当前目录的 config 子目录下
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperV1() {
	cfile := pflag.String("config",
		"config/config.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperReader() {
	cfg := `
db:
  dsn: "root:root@tcp(localhost:13316)/webook"

redis:
  addr: "localhost:6379"
`
	// 读取的类型是 yaml 文件
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}

func initViperWatch() {
	cfile := pflag.String("config",
		"config/config.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	// 实时监听配置变更
	viper.WatchConfig()
	// 只能告诉你文件变了，不能告诉你，文件的哪些内容变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 比较好的设计，它会在 in 里面告诉你变更前的数据，和变更后的数据
		// 更好的设计是，它会直接告诉你差异。
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})
	//viper.SetDefault("db.mysql.dsn",
	//	"root:root@tcp(localhost:3306)/mysql")
	//viper.SetConfigFile("config/dev.yaml")
	//viper.KeyDelimiter("-")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperRemote() {
	err := viper.AddRemoteProvider("etcd3",
		// 通过 ebook 和其他使用 etcd 的区别出来
		"http://127.0.0.1:12379", "/ebook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.WatchRemoteConfig()
	if err != nil {
		panic(err)
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
	})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

//func initServer() *gin.Engine {
//	server := gin.Default()
//
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: config.Config.Redis.Addr,
//	})
//	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
//
//	// 跨域拦截器
//	server.Use(cors.New(cors.Config{
//		//AllowOrigins: []string{"*"},
//		//AllowMethods: []string{"POST", "GET"},
//		AllowHeaders: []string{"Content-Type", "Authorization"},
//		// 不加这个，前端是拿不到 jwt
//		ExposeHeaders: []string{"x-jwt-token"},
//		// 是否允许你带 cookie 之类的东西
//		AllowCredentials: true,
//		AllowOriginFunc: func(origin string) bool {
//			if strings.HasPrefix(origin, "http://localhost") {
//				// 你的开发环境
//				return true
//			}
//			return strings.Contains(origin, "company.com")
//		},
//		MaxAge: 12 * time.Hour,
//	}))
//
//	// 设置 session
//	//// 设置到 cookie
//	//store := cookie.NewStore([]byte("secret"))
//	//// 设置到本地内存
//	store := memstore.NewStore(handler.SessAuthKey, handler.SessEncryptionKey)
//	// 设置到redis
//	//store, err := redisSess.NewStore(16,
//	//	"tcp", "localhost:6379", "",
//	//	handler.SessAuthKey, handler.SessEncryptionKey)
//	//
//	//if err != nil {
//	//	panic(err)
//	//}
//	server.Use(sessions.Sessions("SESS", store))
//
//	// jwt 登录中间件
//	server.Use(middleware.NewJWTLoginMiddlewareBuilder().Build())
//
//	return server
//
//}
//
//func initUser(db *gorm.DB, rcCmd redis.Cmdable) *handler.UserHandler {
//	rc := cache.NewRedisUserCache(rcCmd)
//	codeRc := cache.NewCodeCache(rcCmd)
//	ud := dao.NewGORMUserDAO(db)
//	repo := repository.NewUserRepository(ud, rc)
//	codeRepo := repository.NewCodeRepository(codeRc)
//	svc := service.NewUserService(repo)
//	codeSvc := service.NewCodeService(codeRepo, memory.NewService())
//	u := handler.NewUserHandler(svc, codeSvc)
//	return u
//}
//
//func initDB() *gorm.DB {
//	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
//	if err != nil {
//		// 我只会在初始化过程中 panic
//		// panic 相当于整个 goroutine 结束
//		// 一旦初始化过程出错，应用就不要启动了
//		panic(err)
//	}
//
//	err = dao.InitTables(db)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}
//
//func initRedis() redis.Cmdable {
//	return redis.NewClient(&redis.Options{
//		Addr:     config.Config.Redis.Addr,
//		Password: "", // no password set
//		DB:       0,  // use default DB
//	})
//}
