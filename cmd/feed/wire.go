//go:build wireinject

package main

import (
	"ebook/cmd/feed/events"
	"ebook/cmd/feed/grpc"
	"ebook/cmd/feed/ioc"
	"ebook/cmd/feed/repository"
	"ebook/cmd/feed/repository/cache"
	"ebook/cmd/feed/repository/dao"
	"ebook/cmd/feed/service"
	"ebook/cmd/pkg/appx"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewFeedPushEventDAO,
	dao.NewFeedPullEventDAO,
	cache.NewFeedEventCache,
	repository.NewFeedEventRepo,
)

var thirdProvider = wire.NewSet(
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.InitRedis,
	ioc.InitKafka,
	ioc.InitDB,
	ioc.InitFollowClient,
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		ioc.RegisterHandler,
		service.NewFeedService,
		grpc.NewFeedEventServiceServer,
		events.NewArticleEventConsumer,
		events.NewFeedEventConsumer,
		events.NewFollowerFeedEventConsumer,
		events.NewLikeFeedEventConsumer,
		ioc.InitGRPCServer,
		ioc.NewConsumers,
		wire.Struct(new(appx.App), "GRPCServer", "Consumers"),
	)
	return new(appx.App)
}