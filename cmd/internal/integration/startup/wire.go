//go:build wireinject

package startup

import (
	repository2 "ebook/cmd/interactive/repository"
	cache2 "ebook/cmd/interactive/repository/cache"
	"ebook/cmd/interactive/repository/dao"
	service2 "ebook/cmd/interactive/service"
	events "ebook/cmd/internal/events/article"
	"ebook/cmd/internal/handler"
	ijwt "ebook/cmd/internal/handler/jwt"
	cronJob "ebook/cmd/internal/job"
	"ebook/cmd/internal/repository"
	"ebook/cmd/internal/repository/cache"
	"ebook/cmd/internal/repository/dao/article"
	"ebook/cmd/internal/repository/dao/async_sms"
	"ebook/cmd/internal/repository/dao/job"
	"ebook/cmd/internal/repository/dao/user"
	"ebook/cmd/internal/service"
	"ebook/cmd/internal/service/sms"
	"ebook/cmd/internal/service/sms/async"
	"ebook/cmd/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB,
	InitLogger,
	NewSyncProducer,
	InitKafka,
)
var userSvcProvider = wire.NewSet(
	user.NewGORMUserDAO,
	cache.NewRedisUserCache,
	repository.NewUserRepository,
	service.NewUserService,
)
var articleSvcProvider = wire.NewSet(
	article.NewGORMArticleDAO,
	events.NewKafkaProducer,
	cache.NewRedisArticleCache,
	repository.NewArticleRepository,
	service.NewArticleService,
)
var interactiveSvcProvider = wire.NewSet(
	dao.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
	repository2.NewInteractiveRepository,
	service2.NewInteractiveService,
)
var rankServiceProvider = wire.NewSet(
	service.NewBatchRankingService,
	repository.NewCachedRankingRepository,
	cache.NewRedisRankingCache,
	cache.NewRankingLocalCache,
)

var jobProviderSet = wire.NewSet(
	job.NewGORMJobDAO,
	repository.NewPreemptCronJobRepository,
	service.NewCronJobService,
)

//go:generate wire
func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		userSvcProvider,
		articleSvcProvider,
		interactiveSvcProvider,
		InitInteractiveClient,
		cache.NewCodeCache,
		repository.NewCodeRepository,
		// service 部分
		// 集成测试我们显式指定使用内存实现
		ioc.InitSMSService,

		InitPhantomWechatService,
		service.NewCodeService,

		ioc.NewWechatHandlerConfig,

		handler.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,
		handler.NewUserHandler,
		handler.NewArticleHandler,
		handler.NewObservabilityHandler,

		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}

func InitArticleHandler(dao article.ArticleDAO) *handler.ArticleHandler {
	wire.Build(thirdProvider,
		userSvcProvider,
		interactiveSvcProvider,
		events.NewKafkaProducer,
		cache.NewRedisArticleCache,
		repository.NewArticleRepository,
		InitInteractiveClient,
		service.NewArticleService,
		handler.NewArticleHandler,
	)
	return new(handler.ArticleHandler)
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil, nil)
}

func InitAsyncSmsService(svc sms.Service) *async.Service {
	wire.Build(thirdProvider, repository.NewAsyncSMSRepository,
		async_sms.NewGORMAsyncSmsDAO,
		async.NewService,
	)
	return &async.Service{}
}

func InitRankingService() service.RankingService {
	wire.Build(thirdProvider,
		interactiveSvcProvider,
		articleSvcProvider,
		InitInteractiveClient,
		// 用不上这个 user repo，所以随便搞一个
		wire.InterfaceValue(new(repository.UserRepository),
			&repository.CachedUserRepository{}),
		rankServiceProvider)
	return &service.BatchRankingService{}
}

func InitInteractiveService() service2.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service2.NewInteractiveService(nil, nil)
}

func InitJobScheduler() *cronJob.Scheduler {
	wire.Build(jobProviderSet, thirdProvider, cronJob.NewScheduler)
	return &cronJob.Scheduler{}
}

func InitJwtHandler() ijwt.Handler {
	wire.Build(thirdProvider, ijwt.NewRedisJWTHandler)
	return ijwt.NewRedisJWTHandler(nil)
}
