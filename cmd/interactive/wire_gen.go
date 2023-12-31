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
	db := ioc.InitDB()
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	cmdable := ioc.InitRedis()
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	logger := ioc.InitLogger()
	interactiveRepository := repository.NewInteractiveRepository(interactiveDAO, interactiveCache, logger)
	interactiveService := service.NewInteractiveService(interactiveRepository, logger)
	interactiveServiceServer := grpc.NewInteractiveServiceServer(interactiveService)
	server := ioc.InitGRPCxServer(interactiveServiceServer)
	client := ioc.InitKafka()
	interactiveReadEventConsumer := article.NewInteractiveReadEventConsumer(client, interactiveRepository, logger)
	v := ioc.NewConsumers(interactiveReadEventConsumer)
	app := &App{
		server:    server,
		consumers: v,
	}
	return app
}

// wire.go:

var serviceProvider = wire.NewSet(dao.NewGORMInteractiveDAO, cache.NewRedisInteractiveCache, repository.NewInteractiveRepository, service.NewInteractiveService)

var thirdProvider = wire.NewSet(ioc.InitRedis, ioc.InitDB, ioc.InitLogger, ioc.InitKafka)
