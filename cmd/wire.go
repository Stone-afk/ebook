//go:build wireinject

package main

import (
	"ebook/cmd/internal/handler"
	"ebook/cmd/internal/repository"
	"ebook/cmd/internal/repository/cache"
	"ebook/cmd/internal/repository/dao"
	"ebook/cmd/internal/service"
	"ebook/cmd/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,

		// 初始化 DAO
		dao.NewGORMUserDAO,

		cache.NewRedisUserCache,
		cache.NewCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,

		service.NewUserService,
		service.NewCodeService,

		// 直接基于内存实现
		ioc.InitSMSService,
		handler.NewUserHandler,

		// 你中间件呢？
		// 你注册路由呢？
		// 你这个地方没有用到前面的任何东西
		//gin.Default,
		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)

	return new(gin.Engine)
}