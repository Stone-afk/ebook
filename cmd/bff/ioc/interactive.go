package ioc

import (
	intrv1 "ebook/cmd/api/proto/gen/intr/v1"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitInterActiveClient(etcdClient *etcdv3.Client) intrv1.InteractiveServiceClient {
	type config struct {
		Target string `json:"target"`
		Secure bool   `json:"secure"`
	}
	var cfg config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
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
	return intrv1.NewInteractiveServiceClient(conn)
}
