//go:build wireinject
package main

import (
	"ebook/cmd/pkg/appx"
	"ebook/cmd/tag/grpc"
	"ebook/cmd/tag/ioc"
	"ebook/cmd/tag/repository/cache"
	"ebook/cmd/tag/repository/dao"
	"ebook/cmd/tag/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitKafka,
	ioc.InitProducer,
	ioc.InitEtcdClient,
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		thirdProvider,
		cache.NewRedisTagCache,
		dao.NewGORMTagDAO,
		ioc.InitRepository,
		service.NewTagService,
		grpc.NewTagServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(appx.App), "GRPCServer"),
	)
	return new(appx.App)
}
