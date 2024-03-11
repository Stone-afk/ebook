//go:build wireinject

package main

import (
	"ebook/cmd/pkg/appx"
	"ebook/cmd/reward/events"
	"ebook/cmd/reward/grpc"
	"ebook/cmd/reward/ioc"
	"ebook/cmd/reward/repository"
	"ebook/cmd/reward/repository/cache"
	"ebook/cmd/reward/repository/dao"
	"ebook/cmd/reward/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitEtcdClient,
	ioc.InitRedis,
	ioc.InitKafka)

//go:generate wire
func Init() *appx.App {
	wire.Build(thirdPartySet,
		service.NewWechatNativeRewardService,
		events.NewPaymentEventConsumer,
		ioc.NewConsumers,
		ioc.InitAccountClient,
		ioc.InitGRPCxServer,
		ioc.InitPaymentClient,
		repository.NewRewardRepository,
		cache.NewRewardRedisCache,
		dao.NewRewardGORMDAO,
		grpc.NewRewardServiceServer,
		wire.Struct(new(appx.App), "GRPCServer", "Consumers"),
	)
	return new(appx.App)
}
