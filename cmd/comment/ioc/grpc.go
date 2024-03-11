package ioc

import (
	grpc6 "ebook/cmd/comment/grpc"
	"ebook/cmd/pkg/grpcx/interceptors/logging"
	"ebook/cmd/pkg/grpcx/server"
	"ebook/cmd/pkg/logger"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(commentServer *grpc6.CommentServiceServer,
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
	commentServer.Register(grpcSrv)
	return server.NewGRPCXServer(grpcSrv, etcdClient, l, cfg.Port, "comment", cfg.EtcdTTL)
}
