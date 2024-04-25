package ioc

import (
	smsv1 "ebook/cmd/api/proto/gen/sms/v1"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitSmsRpcClient(etcdClient *etcdv3.Client) smsv1.SmsServiceClient {
	type config struct {
		Target string `json:"target"`
		Secure bool   `json:"secure"`
	}
	var cfg config
	err := viper.UnmarshalKey("grpc.client.sms", &cfg)
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
	cc, err := grpc.Dial(cfg.Target, opts...)
	if err != nil {
		panic(err)
	}
	return smsv1.NewSmsServiceClient(cc)
}
