//go:build wireinject

package main

import (
	"ebook/cmd/cronjob/grpc"
	"ebook/cmd/cronjob/ioc"
	"ebook/cmd/cronjob/repository"
	"ebook/cmd/cronjob/repository/dao"
	"ebook/cmd/cronjob/service"
	"ebook/cmd/pkg/appx"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMJobDAO,
	repository.NewPreemptCronJobRepository,
	service.NewCronJobService)

var thirdProvider = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitEtcdClient,
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewCronJobServiceServer,
		ioc.InitGRPCServer,
		ioc.InitRankingClient,
		ioc.InitExecutors,
		ioc.InitScheduler,
		wire.Struct(new(appx.App), "GRPCServer", "Scheduler"),
	)
	return new(appx.App)
}