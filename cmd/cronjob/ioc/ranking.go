package ioc

import (
	rankingv1 "ebook/cmd/api/proto/gen/ranking/v1"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitRankingClient(etcdClient *etcdv3.Client) rankingv1.RankingServiceClient {
	type config struct {
		Target string `json:"target"`
		Secure bool   `json:"secure"`
	}
	var cfg config
	err := viper.UnmarshalKey("grpc.client.ranking", &cfg)
	if err != nil {
		panic(err)
	}
	rs, err := resolver.NewBuilder(etcdClient)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{grpc.WithResolvers(rs)}
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	conn, err := grpc.Dial(cfg.Target, opts...)
	if err != nil {
		panic(err)
	}
	client := rankingv1.NewRankingServiceClient(conn)
	return client
}
