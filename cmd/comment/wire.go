//go:build wireinject

package main

import (
	grpc2 "ebook/cmd/comment/grpc"
	"ebook/cmd/comment/ioc"
	"ebook/cmd/comment/repository"
	"ebook/cmd/comment/repository/dao"
	"ebook/cmd/comment/service"
	"ebook/cmd/pkg/appx"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewCommentDAO,
	repository.NewCommentRepo,
	service.NewCommentSvc,
	grpc2.NewCommentServer,
)

var thirdProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitEtcdClient,
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		ioc.InitGRPCxServer,
		wire.Struct(new(appx.App), "GRPCServer"),
	)
	return new(appx.App)
}
