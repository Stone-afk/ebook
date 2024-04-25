//go:build wireinject

package code

import (
	"ebook/cmd/code/grpc"
	"ebook/cmd/code/ioc"
	"ebook/cmd/code/repository"
	"ebook/cmd/code/repository/cache"
	"ebook/cmd/code/service"
	"ebook/cmd/pkg/appx"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitEtcdClient,
	ioc.InitLogger,
)

//go:generate wire
func Init() *appx.App {
	wire.Build(
		thirdProvider,
		ioc.InitSmsRpcClient,
		cache.NewCodeCache,
		repository.NewCodeRepository,
		service.NewCodeService,
		grpc.NewCodeServiceServer,
		ioc.InitGRPCServer,
		wire.Struct(new(appx.App), "GRPCServer"),
	)
	return new(appx.App)
}
