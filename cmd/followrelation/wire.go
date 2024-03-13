//go:build wireinject

package main

import (
	grpc7 "ebook/cmd/followrelation/grpc"
	"ebook/cmd/followrelation/ioc"
	"ebook/cmd/followrelation/repository"
	"ebook/cmd/followrelation/repository/cache"
	"ebook/cmd/followrelation/repository/dao"
	"ebook/cmd/followrelation/service"
	"ebook/cmd/pkg/appx"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	cache.NewRedisFollowCache,
	dao.NewGORMFollowRelationDAO,
	repository.NewFollowRelationRepository,
	service.NewFollowRelationService,
	grpc7.NewFollowRelationServiceServer,
)

var thirdProvider = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitEtcdClient,
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		ioc.InitGRPCxServer,
		wire.Struct(new(appx.App), "GRPCServer"),
	)
	return new(appx.App)
}
