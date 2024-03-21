//go:build wireinject

package main

import (
	"ebook/cmd/pkg/appx"
	"ebook/cmd/search/events"
	"ebook/cmd/search/grpc"
	"ebook/cmd/search/ioc"
	"ebook/cmd/search/repository"
	"ebook/cmd/search/repository/dao"
	"ebook/cmd/search/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewUserElasticDAO,
	dao.NewArticleElasticDAO,
	dao.NewAnyElasticDAO,
	dao.NewTagElasticDAO,
	repository.NewUserRepository,
	repository.NewArticleRepository,
	repository.NewAnyRepository,
	repository.NewTagRepository,
	service.NewSyncService,
	service.NewSearchService,
	service.NewArticleSearchService,
	service.NewTagSearchService,
	service.NewUserSearchService,
)

var thirdProvider = wire.NewSet(
	ioc.InitESClient,
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.InitKafka)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewSyncServiceServer,
		grpc.NewSearchServiceServer,
		grpc.NewArticleSearchServiceServer,
		grpc.NewTagSearchServiceServer,
		grpc.NewUserSearchServiceServer,
		events.NewUserConsumer,
		events.NewArticleConsumer,
		events.NewSyncDataEventConsumer,
		events.NewTagConsumer,
		ioc.InitGRPCxServer,
		ioc.NewConsumers,
		wire.Struct(new(appx.App), "GRPCServer", "Consumers"),
	)
	return new(appx.App)
}
