package main

import (
	"ebook/cmd/pkg/ginx"
	"ebook/cmd/pkg/grpcx/server"
	"ebook/cmd/pkg/saramax"
)

type App struct {
	// 在这里，所有需要 main 函数控制启动、关闭的，都会在这里有一个
	// 核心就是为了控制生命周期
	server    *server.Server
	webAdmin  *ginx.Server
	consumers []saramax.Consumer
}
