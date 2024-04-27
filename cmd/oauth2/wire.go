//go:build wireinject

package main

import (
	"ebook/cmd/oauth2/grpc"
	"ebook/cmd/oauth2/ioc"
	"ebook/cmd/pkg/appx"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitEtcdClient,
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		thirdProvider,
		ioc.InitPrometheus,
		grpc.NewOauth2ServiceServer,
		ioc.InitGRPCServer,
		wire.Struct(new(appx.App), "GRPCServer"),
	)
	return new(appx.App)
}
