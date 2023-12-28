//go:build wireinject

package main

import (
	repository2 "ebook/cmd/interactive/repository"
	cache2 "ebook/cmd/interactive/repository/cache"
	"ebook/cmd/interactive/repository/dao"
	service2 "ebook/cmd/interactive/service"
	events "ebook/cmd/internal/events/article"
	"ebook/cmd/internal/handler"
	ijwt "ebook/cmd/internal/handler/jwt"
	"ebook/cmd/internal/repository"
	"ebook/cmd/internal/repository/cache"
	"ebook/cmd/internal/repository/dao/article"
	"ebook/cmd/internal/repository/dao/job"
	"ebook/cmd/internal/repository/dao/user"
	"ebook/cmd/internal/service"
	"ebook/cmd/ioc"
	"github.com/google/wire"
)

var rankServiceProvider = wire.NewSet(
	service.NewBatchRankingService,
	repository.NewCachedRankingRepository,
	cache.NewRedisRankingCache,
	cache.NewRankingLocalCache,
)

var interactiveServiceProvider = wire.NewSet(
	dao.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
	repository2.NewInteractiveRepository,
	service2.NewInteractiveService,
)

var articleServiceProvider = wire.NewSet(
	article.NewGORMArticleDAO,
	cache.NewRedisArticleCache,
	repository.NewArticleRepository,
	service.NewArticleService,
)

var userServiceProvider = wire.NewSet(
	// 初始化 DAO
	user.NewGORMUserDAO,
	cache.NewRedisUserCache,
	repository.NewUserRepository,
	service.NewUserService,
)

var codeServiceProvider = wire.NewSet(
	cache.NewCodeCache,
	repository.NewCodeRepository,
	ioc.InitSMSService,
	service.NewCodeService,
)

var wechatServiceProvider = wire.NewSet(
	ioc.InitWechatService,
	ioc.NewWechatHandlerConfig,
)

var jobSvcProvider = wire.NewSet(
	job.NewGORMJobDAO,
	repository.NewPreemptCronJobRepository,
	service.NewCronJobService,
)

//go:generate wire
func InitApp() *App {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,
		ioc.InitRLockClient,

		ioc.InitJobs,
		ioc.InitRankingJob,

		repository.NewHistoryRecordRepository,
		// events
		events.NewHistoryReadEventConsumer,
		events.NewKafkaProducer,

		rankServiceProvider,
		interactiveServiceProvider,
		articleServiceProvider,
		userServiceProvider,
		codeServiceProvider,
		wechatServiceProvider,
		ioc.InitInteractiveGRPCClient,

		handler.NewOAuth2WechatHandler,
		handler.NewUserHandler,
		handler.NewArticleHandler,
		handler.NewObservabilityHandler,
		ijwt.NewRedisJWTHandler,
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
