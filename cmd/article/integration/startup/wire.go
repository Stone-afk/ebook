//go:build wireinject
package startup

import (
	"ebook/cmd/article/events"
	"ebook/cmd/article/repository"
	"ebook/cmd/article/repository/cache"
	"ebook/cmd/article/repository/dao"
	"ebook/cmd/article/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB,
	InitLogger,
	NewSyncProducer,
	InitKafka,
)

var articleSvcProvider = wire.NewSet(
	dao.NewGORMArticleDAO,
	events.NewKafkaProducer,
	cache.NewRedisArticleCache,
	repository.NewArticleRepository,
	service.NewArticleService,
)

type TestArticleService struct {
	service.ArticleService
}

//go:generate wire
func InitArticleService() service.ArticleService {
	wire.Build(thirdProvider, articleSvcProvider)
	// 随便返回一个
	return new(TestArticleService)
}
