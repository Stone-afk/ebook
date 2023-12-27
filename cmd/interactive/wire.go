//go:build wireinject

package interactive

import (
	events "ebook/cmd/interactive/events/article"
	"ebook/cmd/interactive/grpc"
	"ebook/cmd/interactive/ioc"
	"ebook/cmd/interactive/repository"
	"ebook/cmd/interactive/repository/cache"
	"ebook/cmd/interactive/repository/dao"
	"ebook/cmd/interactive/service"
	"github.com/google/wire"
)

var serviceProvider = wire.NewSet(
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
	repository.NewInteractiveRepository,
	service.NewInteractiveService)

var thirdProvider = wire.NewSet(
	ioc.InitRedis, ioc.InitDB,
	ioc.InitLogger,
	ioc.InitKafka,
)

//go:generate wire
func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProvider,
		events.NewInteractiveReadEventConsumer,
		grpc.NewInteractiveServiceServer,
		ioc.NewConsumers,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
