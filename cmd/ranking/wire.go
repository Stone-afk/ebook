//go:build wireinject

package main

import (
	"ebook/cmd/pkg/appx"
	"ebook/cmd/ranking/grpc"
	"ebook/cmd/ranking/ioc"
	"ebook/cmd/ranking/repository"
	"ebook/cmd/ranking/repository/cache"
	"ebook/cmd/ranking/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	cache.NewRankingLocalCache,
	cache.NewRedisRankingCache,
	repository.NewCachedRankingRepository,
	service.NewBatchRankingService,
)

var thirdProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitEtcdClient,
	ioc.InitRedis,
	ioc.InitInterActiveRpcClient,
	ioc.InitArticleRpcClient,
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewRankingServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(appx.App), "GRPCServer"),
	)
	return new(appx.App)
}
