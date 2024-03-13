//go:build wireinject
package startup

import (
	grpc2 "ebook/cmd/followrelation/grpc"
	"ebook/cmd/followrelation/repository"
	"ebook/cmd/followrelation/repository/cache"
	"ebook/cmd/followrelation/repository/dao"
	"ebook/cmd/followrelation/service"
	"github.com/google/wire"
)

//go:generate wire
func InitServer() *grpc2.FollowServiceServer {
	wire.Build(
		InitRedis,
		InitLogger,
		InitTestDB,
		dao.NewGORMFollowRelationDAO,
		cache.NewRedisFollowCache,
		repository.NewFollowRelationRepository,
		service.NewFollowRelationService,
		grpc2.NewFollowRelationServiceServer,
	)
	return new(grpc2.FollowServiceServer)
}