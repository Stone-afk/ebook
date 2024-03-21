//go:build wireinject

package main

import (
	"ebook/cmd/pkg/appx"
	"ebook/cmd/user/events"
	"ebook/cmd/user/grpc"
	"ebook/cmd/user/ioc"
	"ebook/cmd/user/repository"
	"ebook/cmd/user/repository/cache"
	"ebook/cmd/user/repository/dao"
	"ebook/cmd/user/service"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitEtcdClient,
	ioc.InitKafka,
	ioc.NewSyncProducer,
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		thirdProvider,
		events.NewSaramaSyncProducer,
		cache.NewRedisUserCache,
		dao.NewGORMUserDAO,
		repository.NewUserRepository,
		service.NewUserService,
		grpc.NewUserServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(appx.App), "GRPCServer"),
	)
	return new(appx.App)
}
