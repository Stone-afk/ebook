//go:build wireinject

package startup

import (
	"ebook/cmd/search/grpc"
	"ebook/cmd/search/ioc"
	"ebook/cmd/search/repository"
	"ebook/cmd/search/repository/dao"
	"ebook/cmd/search/service"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewUserElasticDAO,
	dao.NewArticleElasticDAO,
	dao.NewTagElasticDAO,
	dao.NewAnyElasticDAO,
	repository.NewUserRepository,
	repository.NewAnyRepository,
	repository.NewArticleRepository,
	repository.NewTagRepository,
	service.NewSyncService,
	service.NewSearchService,
	service.NewArticleSearchService,
	service.NewTagSearchService,
	service.NewUserSearchService,
)

var thirdProvider = wire.NewSet(
	InitESClient,
	ioc.InitLogger)

//go:generate wire
func InitTagSearchServer() *grpc.TagSearchServiceServer {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewTagSearchServiceServer,
	)
	return new(grpc.TagSearchServiceServer)
}

func InitUserSearchServer() *grpc.UserSearchServiceServer {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewUserSearchServiceServer,
	)
	return new(grpc.UserSearchServiceServer)
}

func InitArticleSearchServer() *grpc.ArticleSearchServiceServer {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewArticleSearchServiceServer,
	)
	return new(grpc.ArticleSearchServiceServer)
}

func InitSearchServer() *grpc.SearchServiceServer {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewSearchServiceServer,
	)
	return new(grpc.SearchServiceServer)
}

func InitSyncServer() *grpc.SyncServiceServer {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewSyncServiceServer,
	)
	return new(grpc.SyncServiceServer)
}
