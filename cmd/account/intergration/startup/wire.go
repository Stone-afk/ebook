//go:build wireinject

package startup

import (
	"ebook/cmd/account/grpc"
	"ebook/cmd/account/repository"
	"ebook/cmd/account/repository/dao"
	"ebook/cmd/account/service"
	"github.com/google/wire"
)

//go:generate wire
func InitAccountService() *grpc.AccountServiceServer {
	wire.Build(InitTestDB,
		dao.NewCreditGORMDAO,
		repository.NewAccountRepository,
		service.NewAccountService,
		grpc.NewAccountServiceServer)
	return new(grpc.AccountServiceServer)
}
