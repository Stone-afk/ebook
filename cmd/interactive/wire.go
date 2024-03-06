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
	ioc.InitDST,
	ioc.InitSRC,
	ioc.InitBizDB,
	ioc.InitDoubleWritePool,
	ioc.InitRedis,
	ioc.InitLogger,
	// 暂时不理会 consumer 怎么启动
	ioc.InitSyncProducer,
	ioc.InitKafka,
	ioc.InitEtcdClient,
)

var migratorProvider = wire.NewSet(
	ioc.InitMigratorWeb,
	ioc.InitFixDataConsumer,
	ioc.InitMigradatorProducer)

//go:generate wire
func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProvider,
		migratorProvider,
		events.NewInteractiveReadEventConsumer,
		grpc.NewInteractiveServiceServer,
		ioc.NewConsumers,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
