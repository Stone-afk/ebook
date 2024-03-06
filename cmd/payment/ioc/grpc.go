package ioc

import (
	grpc2 "ebook/cmd/payment/grpc"
	"ebook/cmd/pkg/grpcx/interceptors/logging"
	"ebook/cmd/pkg/grpcx/server"
	"ebook/cmd/pkg/logger"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCServer(l logger.Logger, etcdClient *etcdv3.Client, weServer *grpc2.WechatServiceServer) *server.Server {
	type Config struct {
		Port    int   `yaml:"port"`
		EtcdTTL int64 `yaml:"etcdTTL"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	grpcSvc := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.NewLoggerInterceptorBuilder(l).Build(),
	))
	weServer.Register(grpcSvc)
	return server.NewGRPCXServer(grpcSvc, etcdClient, l, cfg.Port, "payment", cfg.EtcdTTL)
}
