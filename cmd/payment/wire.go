//go:build wireinject

package main

import (
	"ebook/cmd/payment/grpc"
	"ebook/cmd/payment/handler"
	"ebook/cmd/payment/ioc"
	"ebook/cmd/payment/repository"
	"ebook/cmd/payment/repository/dao"
	"ebook/cmd/pkg/appx"
	"github.com/google/wire"
)

//go:generate wire
func InitApp() *appx.App {
	wire.Build(
		ioc.InitEtcdClient,
		ioc.InitKafka,
		ioc.InitProducer,
		ioc.InitWechatClient,
		dao.NewPaymentGORMDAO,
		ioc.InitDB,
		repository.NewPaymentRepository,
		grpc.NewWechatServiceServer,
		ioc.InitWechatNativeService,
		ioc.InitWechatConfig,
		ioc.InitWechatNotifyHandler,
		ioc.InitGRPCServer,
		handler.NewWechatHandler,
		ioc.InitGinServer,
		ioc.InitLogger,
		wire.Struct(new(appx.App), "WebServer", "GRPCServer"))
	return new(appx.App)
}
