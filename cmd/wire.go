//go:build wireinject

package main

import (
	events "ebook/cmd/internal/events/article"
	"ebook/cmd/internal/handler"
	ijwt "ebook/cmd/internal/handler/jwt"
	"ebook/cmd/internal/repository"
	"ebook/cmd/internal/repository/cache"
	"ebook/cmd/internal/repository/dao/article"
	"ebook/cmd/internal/repository/dao/interactive"
	"ebook/cmd/internal/repository/dao/user"
	"ebook/cmd/internal/service"
	"ebook/cmd/ioc"
	"github.com/google/wire"
)

func InitApp() *App {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,

		// consumer
		events.NewInteractiveReadEventBatchConsumer,
		events.NewKafkaProducer,

		// 初始化 DAO
		user.NewGORMUserDAO,
		article.NewGORMArticleDAO,
		interactive.NewGORMInteractiveDAO,

		cache.NewRedisInteractiveCache,
		cache.NewRedisUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository.NewInteractiveRepository,
		repository.NewArticleRepository,

		service.NewUserService,
		service.NewCodeService,
		ioc.InitSMSService,
		ioc.InitWechatService,
		service.NewArticleService,
		service.NewInteractiveService,

		ioc.NewWechatHandlerConfig,

		handler.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,
		handler.NewUserHandler,
		handler.NewArticleHandler,
		// 你中间件呢？
		// 你注册路由呢？
		// 你这个地方没有用到前面的任何东西
		//gin.Default,
		ioc.InitMiddlewares,
		ioc.InitWebServer,
		// 组装我这个结构体的所有字段
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
