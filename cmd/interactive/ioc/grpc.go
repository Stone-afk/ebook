package ioc

import (
	grpc2 "ebook/cmd/interactive/grpc"
	"ebook/cmd/pkg/grpcx/server"
	"ebook/cmd/pkg/logger"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(l logger.Logger, etcdClient *etcdv3.Client, intrServer *grpc2.InteractiveServiceServer) *server.Server {
	type Config struct {
		Port    int   `yaml:"port"`
		EtcdTTL int64 `yaml:"etcdTTL"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	// master 分支
	//err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}
	grpcSvc := grpc.NewServer()
	intrServer.Registry(grpcSvc)
	return server.NewGRPCXServer(grpcSvc, etcdClient, l, cfg.Port, "interactive", cfg.EtcdTTL)
}

//func InitZeroServer(intrServer *grpc2.InteractiveServiceServer) *server.ZeroServer {
//	type Config struct {
//		Port      int      `yaml:"port"`
//		EtcdAddrs []string `yaml:"etcdAddrs"`
//	}
//	var cfg Config
//	err := viper.UnmarshalKey("grpc.server", &cfg)
//	// master 分支
//	//err := viper.UnmarshalKey("grpc", &cfg)
//	if err != nil {
//		panic(err)
//	}
//
//	return &server.ZeroServer{
//		Port:      cfg.Port,
//		EtcdAddrs: cfg.EtcdAddrs,
//		Name:      "interactive",
//		Register:  intrServer,
//	}
//}
//
//func InitKratosServer(intrServer *grpc2.InteractiveServiceServer) *server.KratosServer {
//	type Config struct {
//		Port      int      `yaml:"port"`
//		EtcdAddrs []string `yaml:"etcdAddrs"`
//	}
//	var cfg Config
//	err := viper.UnmarshalKey("grpc.server", &cfg)
//	// master 分支
//	//err := viper.UnmarshalKey("grpc", &cfg)
//	if err != nil {
//		panic(err)
//	}
//
//	return &server.KratosServer{
//		Port:      cfg.Port,
//		EtcdAddrs: cfg.EtcdAddrs,
//		Name:      "interactive",
//		Register:  intrServer,
//	}
//}
