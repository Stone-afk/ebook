package ioc

import (
	"ebook/cmd/payment/handler"
	"ebook/cmd/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

func InitGinServer(hdl *handler.WechatHandler) *ginx.Server {
	engine := gin.Default()
	hdl.RegisterRoutes(engine)
	addr := viper.GetString("http.addr")
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "ebook",
		Subsystem: "payment",
		Name:      "http",
	})
	return &ginx.Server{
		Engine: engine,
		Addr:   addr,
	}
}
