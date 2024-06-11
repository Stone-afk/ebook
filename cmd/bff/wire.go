//go:build wireinject

package main

import (
	"ebook/cmd/bff/handler"
	"ebook/cmd/bff/handler/jwt"
	"ebook/cmd/bff/ioc"
	"ebook/cmd/pkg/appx"
	"github.com/google/wire"
)

//go:generate wire
func InitApp() *appx.App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitRedis,
		ioc.InitEtcdClient,

		handler.NewArticleHandler,
		handler.NewUserHandler,
		handler.NewRewardHandler,
		handler.NewOAuth2WechatHandler,
		handler.NewObservabilityHandler,
		jwt.NewRedisJWTHandler,

		ioc.InitMiddlewares,
		ioc.InitUserClient,
		ioc.InitInterActiveClient,
		ioc.InitOauth2ServiceClient,
		ioc.NewWechatHandlerConfig,
		ioc.InitRewardClient,
		ioc.InitCodeClient,
		ioc.InitArticleClient,
		ioc.InitGinServer,
		wire.Struct(new(appx.App), "WebServer"),
	)
	return new(appx.App)
}
