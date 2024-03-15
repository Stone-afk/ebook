//go:build wireinject
package startup

import (
	"ebook/cmd/tag/grpc"
	"ebook/cmd/tag/repository/cache"
	"ebook/cmd/tag/repository/dao"
	"ebook/cmd/tag/service"
	"github.com/google/wire"
)

//go:generate wire
func InitGRPCService() *grpc.TagServiceServer {
	wire.Build(InitTestDB, InitRedis,
		InitLogger,
		dao.NewGORMTagDAO,
		InitRepository,
		cache.NewRedisTagCache,
		service.NewTagService,
		grpc.NewTagServiceServer,
	)
	return new(grpc.TagServiceServer)
}
