//go:build wireinject

package startup

import (
	grpc2 "ebook/cmd/comment/grpc"
	"ebook/cmd/comment/repository"
	"ebook/cmd/comment/repository/dao"
	"ebook/cmd/comment/service"
	"ebook/cmd/pkg/logger/nop"
	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewCommentDAO,
	repository.NewCommentRepo,
	service.NewCommentSvc,
	grpc2.NewCommentServer,
)

var thirdProvider = wire.NewSet(
	nop.NewNopLogger,
	InitTestDB,
)

//go:generate wire
func InitGRPCServer() *grpc2.CommentServiceServer {
	wire.Build(thirdProvider, serviceProviderSet)
	return new(grpc2.CommentServiceServer)
}
