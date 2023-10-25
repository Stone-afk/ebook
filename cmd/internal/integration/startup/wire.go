//go:build wireinject

package startup

import (
	"ebook/cmd/internal/handler"
	ijwt "ebook/cmd/internal/handler/jwt"
	"ebook/cmd/internal/repository"
	"ebook/cmd/internal/repository/cache"
	"ebook/cmd/internal/repository/dao/article"
	"ebook/cmd/internal/repository/dao/user"
	"ebook/cmd/internal/service"
	"ebook/cmd/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLogger)
var userSvcProvider = wire.NewSet(
	user.NewGORMUserDAO,
	cache.NewRedisUserCache,
	repository.NewUserRepository,
	service.NewUserService)
var articleSvcProvider = wire.NewSet(
	article.NewGORMArticleDAO,
	repository.NewArticleRepository,
	service.NewArticleService,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		userSvcProvider,
		articleSvcProvider,

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

		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)
	// 随便返回一个
	return gin.Default()
}

func InitArticleHandler(dao article.ArticleDAO) *handler.ArticleHandler {
	wire.Build(thirdProvider,
		//userSvcProvider,
		repository.NewArticleRepository,
		service.NewArticleService,
		handler.NewArticleHandler,
	)
	return new(handler.ArticleHandler)
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil, nil)
}

func InitJwtHandler() ijwt.Handler {
	wire.Build(thirdProvider, ijwt.NewRedisJWTHandler)
	return ijwt.NewRedisJWTHandler(nil)
}
