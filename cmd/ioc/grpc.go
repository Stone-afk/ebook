package ioc

import (
	intrv1 "ebook/cmd/api/proto/gen/intr/v1"
	"ebook/cmd/interactive/service"
	"ebook/cmd/internal/handler/client/interactive"
	"ebook/cmd/pkg/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitEtcd() *etcdv3.Client {
	var cfg etcdv3.Config
	err := viper.UnmarshalKey("etcd", &cfg)
	if err != nil {
		panic(err)
	}
	cli, err := etcdv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return cli
}

// InitInteractiveGRPCClient 真正的 gRPC 的客户端
func InitInteractiveGRPCClient(svc service.InteractiveService, l logger.Logger) intrv1.InteractiveServiceClient {
	panic("")
}

// InitInteractiveGRPCClientV2 这个是流量控制的客户端
func InitInteractiveGRPCClientV2(svc service.InteractiveService, l logger.Logger) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr      string
		Secure    bool
		Threshold int32
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	if cfg.Secure {
		// 上面，要去加载证书之类的东西
		// 启用 HTTPS
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}
	remoteClient := intrv1.NewInteractiveServiceClient(cc)
	localClient := interactive.NewServiceAdapter(svc)
	res := interactive.NewGreyScaleServiceClient(remoteClient, localClient, cfg.Threshold)
	// 在这里监听
	viper.OnConfigChange(func(in fsnotify.Event) {
		cfg = Config{}
		err1 := viper.UnmarshalKey("grpc.intr", cfg)
		if err1 != nil {
			l.Error("重新加载grpc.intr的配置失败", logger.Error(err1))
			return
		}
		// 这边更新 Threshold
		res.UpdateThreshold(cfg.Threshold)
	})
	return res
}
