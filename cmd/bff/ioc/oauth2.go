package ioc

import (
	oauth2v1 "ebook/cmd/api/proto/gen/oauth2/v1"
	"ebook/cmd/bff/handler"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitOauth2ServiceClient(ecli *clientv3.Client) oauth2v1.Oauth2ServiceClient {
	type Config struct {
		Target string `json:"target"`
		Secure bool   `json:"secure"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.code", &cfg)
	if err != nil {
		panic(err)
	}
	rs, err := resolver.NewBuilder(ecli)
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
	return oauth2v1.NewOauth2ServiceClient(cc)
}

func NewWechatHandlerConfig() handler.WechatHandlerConfig {
	return handler.WechatHandlerConfig{
		Secure: false,
	}
}
