//go:build wireinject

package main

import (
	"ebook/cmd/article/events"
	"ebook/cmd/article/grpc"
	"ebook/cmd/article/ioc"
	"ebook/cmd/article/repository"
	"ebook/cmd/article/repository/cache"
	"ebook/cmd/article/repository/dao"
	"ebook/cmd/article/service"
	"ebook/cmd/pkg/appx"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitUserClient,
	ioc.InitEtcdClient,
	ioc.InitDB,
	ioc.InitKafka,
	ioc.NewSyncProducer,
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		thirdProvider,
		events.NewKafkaProducer,
		events.NewMySQLBinlogConsumer,
		events.NewSaramaSyncProducer,
		cache.NewRedisArticleCache,
		dao.NewGORMArticleDAO,
		repository.NewCachedArticleRepository,
		repository.NewCachedArticleRepo,
		service.NewArticleService,
		grpc.NewArticleServiceServer,
		ioc.InitGRPCxServer,
		ioc.NewConsumers,
		wire.Struct(new(appx.App), "GRPCServer", "Consumers"),
	)
	return new(appx.App)
}
