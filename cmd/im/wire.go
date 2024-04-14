//go:build wireinject

package main

import (
	"ebook/cmd/im/events"
	"ebook/cmd/im/ioc"
	"ebook/cmd/im/service"
	"ebook/cmd/pkg/appx"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitKafka,
	ioc.InitLogger,
	ioc.GetSecret,
	ioc.GetBaseHost,
	ioc.InitDefaultHttpClient,
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		thirdProvider,
		events.NewSyncUserConsumer,
		service.NewRESTUserService,
		ioc.NewConsumers,
		wire.Struct(new(appx.App), "Consumers"),
	)
	return new(appx.App)
}
