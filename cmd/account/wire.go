//go:build wireinject
package main

import (
	"ebook/cmd/account/grpc"
	"ebook/cmd/account/ioc"
	"ebook/cmd/account/repository"
	"ebook/cmd/account/repository/dao"
	"ebook/cmd/account/service"
	"ebook/cmd/pkg/appx"
	"github.com/google/wire"
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		ioc.InitDB,
		ioc.InitLogger,
		ioc.InitEtcdClient,
		ioc.InitGRPCxServer,
		dao.NewCreditGORMDAO,
		repository.NewAccountRepository,
		service.NewAccountService,
		grpc.NewAccountServiceServer,
		wire.Struct(new(appx.App), "GRPCServer"))
	return new(appx.App)
}
