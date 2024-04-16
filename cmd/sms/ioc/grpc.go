package ioc

import (
	"ebook/cmd/pkg/grpcx/server"
	"ebook/cmd/pkg/logger"
	grpc2 "ebook/cmd/sms/grpc"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(reward *grpc2.SmsServiceServer,
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
	grpcSrv := grpc.NewServer()
	reward.Register(grpcSrv)
	return server.NewGRPCXServer(grpcSrv, etcdClient, l, cfg.Port, "sms", cfg.EtcdTTL)
}
