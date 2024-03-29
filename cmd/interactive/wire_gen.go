// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package interactive

import (
	"ebook/cmd/interactive/events/article"
	"ebook/cmd/interactive/grpc"
	"ebook/cmd/interactive/ioc"
	"ebook/cmd/interactive/repository"
	"ebook/cmd/interactive/repository/cache"
	"ebook/cmd/interactive/repository/dao"
	"ebook/cmd/interactive/service"
	"github.com/google/wire"
)

// Injectors from wire.go:

//go:generate wire
func Init() *App {
	logger := ioc.InitLogger()
	client := ioc.InitEtcdClient()
	srcDB := ioc.InitSRC(logger)
	dstDB := ioc.InitDST(logger)
	doubleWritePool := ioc.InitDoubleWritePool(srcDB, dstDB)
	db := ioc.InitBizDB(doubleWritePool)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	cmdable := ioc.InitRedis()
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	interactiveRepository := repository.NewInteractiveRepository(interactiveDAO, interactiveCache, logger)
	interactiveService := service.NewInteractiveService(interactiveRepository, logger)
	interactiveServiceServer := grpc.NewInteractiveServiceServer(interactiveService)
	server := ioc.InitGRPCxServer(logger, client, interactiveServiceServer)
	saramaClient := ioc.InitKafka()
	syncProducer := ioc.InitSyncProducer(saramaClient)
	producer := ioc.InitMigradatorProducer(syncProducer)
	ginxServer := ioc.InitMigratorWeb(logger, srcDB, dstDB, doubleWritePool, producer)
	interactiveReadEventConsumer := article.NewInteractiveReadEventConsumer(saramaClient, interactiveRepository, logger)
	consumer := ioc.InitFixDataConsumer(logger, srcDB, dstDB, saramaClient)
	v := ioc.NewConsumers(interactiveReadEventConsumer, consumer)
	app := &App{
		server:    server,
		webAdmin:  ginxServer,
		consumers: v,
	}
	return app
}

// wire.go:

var serviceProvider = wire.NewSet(dao.NewGORMInteractiveDAO, cache.NewRedisInteractiveCache, repository.NewInteractiveRepository, service.NewInteractiveService)

var thirdProvider = wire.NewSet(ioc.InitDST, ioc.InitSRC, ioc.InitBizDB, ioc.InitDoubleWritePool, ioc.InitRedis, ioc.InitLogger, ioc.InitSyncProducer, ioc.InitKafka, ioc.InitEtcdClient)

var migratorProvider = wire.NewSet(ioc.InitMigratorWeb, ioc.InitFixDataConsumer, ioc.InitMigradatorProducer)
