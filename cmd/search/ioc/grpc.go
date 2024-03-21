package ioc

import (
	"ebook/cmd/pkg/grpcx/interceptors/logging"
	"ebook/cmd/pkg/grpcx/server"
	"ebook/cmd/pkg/logger"
	grpc9 "ebook/cmd/search/grpc"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(syncServiceServer *grpc9.SyncServiceServer,
	searchServiceServer *grpc9.SearchServiceServer,
	userSearchServiceServer *grpc9.UserSearchServiceServer,
	articleSearchServiceServer *grpc9.ArticleSearchServiceServer,
	tagSearchServiceServer *grpc9.TagSearchServiceServer,
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
	syncServiceServer.Register(grpcSrv)
	searchServiceServer.Register(grpcSrv)
	userSearchServiceServer.Register(grpcSrv)
	articleSearchServiceServer.Register(grpcSrv)
	tagSearchServiceServer.Register(grpcSrv)
	return server.NewGRPCXServer(grpcSrv, etcdClient, l, cfg.Port, "search", cfg.EtcdTTL)
}
