//go:build wireinject

package sms

import (
	"ebook/cmd/pkg/appx"
	"ebook/cmd/sms/grpc"
	"ebook/cmd/sms/ioc"
	"github.com/google/wire"
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitEtcdClient,
		ioc.InitSmsTencentService,
		grpc.NewSmsServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(appx.App), "GRPCServer"),
	)
	return new(appx.App)
}
