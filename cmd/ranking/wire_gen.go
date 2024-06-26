// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"ebook/cmd/pkg/appx"
	"ebook/cmd/ranking/grpc"
	"ebook/cmd/ranking/ioc"
	"ebook/cmd/ranking/repository"
	"ebook/cmd/ranking/repository/cache"
	"ebook/cmd/ranking/service"
	"github.com/google/wire"
)

// Injectors from wire.go:

//go:generate wire
func Init() *appx.App {
	client := ioc.InitEtcdClient()
	interactiveServiceClient := ioc.InitInterActiveRpcClient(client)
	articleServiceClient := ioc.InitArticleRpcClient(client)
	cmdable := ioc.InitRedis()
	redisRankingCache := cache.NewRedisRankingCache(cmdable)
	rankingLocalCache := cache.NewRankingLocalCache()
	rankingRepository := repository.NewCachedRankingRepository(redisRankingCache, rankingLocalCache)
	rankingService := service.NewBatchRankingService(interactiveServiceClient, articleServiceClient, rankingRepository)
	rankingServiceServer := grpc.NewRankingServiceServer(rankingService)
	logger := ioc.InitLogger()
	server := ioc.InitGRPCxServer(rankingServiceServer, client, logger)
	app := &appx.App{
		GRPCServer: server,
	}
	return app
}

// wire.go:

var serviceProviderSet = wire.NewSet(cache.NewRankingLocalCache, cache.NewRedisRankingCache, repository.NewCachedRankingRepository, service.NewBatchRankingService)

var thirdProvider = wire.NewSet(ioc.InitLogger, ioc.InitEtcdClient, ioc.InitRedis, ioc.InitInterActiveRpcClient, ioc.InitArticleRpcClient)
