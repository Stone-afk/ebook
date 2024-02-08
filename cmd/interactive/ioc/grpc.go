package ioc

import (
	grpc2 "ebook/cmd/interactive/grpc"
	"ebook/cmd/pkg/grpcx"
	"ebook/cmd/pkg/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(l logger.Logger, intrServer *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Port      int      `yaml:"port"`
		EtcdAddrs []string `yaml:"etcdAddrs"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	// master 分支
	//err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	intrServer.Registry(server)

	return &grpcx.Server{
		Server:    server,
		Port:      cfg.Port,
		EtcdAddrs: cfg.EtcdAddrs,
		Name:      "interactive",
		L:         l,
	}
}

func InitZeroServer(intrServer *grpc2.InteractiveServiceServer) *grpcx.ZeroServer {
	type Config struct {
		Port      int      `yaml:"port"`
		EtcdAddrs []string `yaml:"etcdAddrs"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	// master 分支
	//err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}

	return &grpcx.ZeroServer{
		Port:      cfg.Port,
		EtcdAddrs: cfg.EtcdAddrs,
		Name:      "interactive",
		Register:  intrServer,
	}
}

func InitKratosServer(intrServer *grpc2.InteractiveServiceServer) *grpcx.KratosServer {
	type Config struct {
		Port      int      `yaml:"port"`
		EtcdAddrs []string `yaml:"etcdAddrs"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	// master 分支
	//err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}

	return &grpcx.KratosServer{
		Port:      cfg.Port,
		EtcdAddrs: cfg.EtcdAddrs,
		Name:      "interactive",
		Register:  intrServer,
	}
}
