package ioc

import (
	grpc5 "ebook/cmd/article/grpc"
	"ebook/cmd/pkg/grpcx/interceptors/logging"
	"ebook/cmd/pkg/grpcx/server"
	"ebook/cmd/pkg/logger"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(articleServer *grpc5.ArticleServiceServer,
	etcdClient *clientv3.Client,
	l logger.Logger) *server.Server {
	type Config struct {
		Port    int   `yaml:"port"`
		EtcdTTL int64 `yaml:"etcdTTL"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	grpcSrv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.NewLoggerInterceptorBuilder(l).Build(),
	))
	articleServer.Register(grpcSrv)
	return server.NewGRPCXServer(grpcSrv, etcdClient, l, cfg.Port, "article", cfg.EtcdTTL)
}
