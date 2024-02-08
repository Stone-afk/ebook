package ioc

import (
	intrv1 "ebook/cmd/api/proto/gen/intr/v1"
	"github.com/spf13/viper"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitEtcd() *clientv3.Client {
	var cfg clientv3.Config
	err := viper.UnmarshalKey("etcd", &cfg)
	if err != nil {
		panic(err)
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return cli
}

// InitInteractiveGRPCClient 真正的 gRPC 的客户端
func InitInteractiveGRPCClient(client *clientv3.Client) intrv1.InteractiveServiceClient {
	type Config struct {
		Secure bool
		Name   string
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	// build 注册中心
	bd, err := resolver.NewBuilder(client)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{grpc.WithResolvers(bd)}
	if cfg.Secure {
		// 上面，要去加载你的证书之类的东西
		// 启用 HTTPS
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	// 连接注册中心
	// 这个地方没填对，它也不会报错
	cc, err := grpc.Dial("etcd:///service/"+cfg.Name, opts...)
	if err != nil {
		panic(err)
	}
	return intrv1.NewInteractiveServiceClient(cc)
}

// InitInteractiveZeroClient 真正的 gRPC 的客户端
func InitInteractiveZeroClient(client *clientv3.Client) intrv1.InteractiveServiceClient {
	type Config struct {
		EtcdAddrs []string `yaml:"etcdAddrs"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	c := zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: cfg.EtcdAddrs,
			Key:   "interactive",
		},
	}
	zClient := zrpc.MustNewClient(c)
	cc := zClient.Conn()
	return intrv1.NewInteractiveServiceClient(cc)
}

//// InitInteractiveThresholdGRPCClient 这个是流量控制的客户端
//func InitInteractiveThresholdGRPCClient(svc service.InteractiveService, l logger.Logger) intrv1.InteractiveServiceClient {
//	type Config struct {
//		Addr      string
//		Secure    bool
//		Threshold int32
//	}
//	var cfg Config
//	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
//	if err != nil {
//		panic(err)
//	}
//	var opts []grpc.DialOption
//	if cfg.Secure {
//		// 上面，要去加载证书之类的东西
//		// 启用 HTTPS
//	} else {
//		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
//	}
//	cc, err := grpc.Dial(cfg.Addr, opts...)
//	if err != nil {
//		panic(err)
//	}
//	remoteClient := intrv1.NewInteractiveServiceClient(cc)
//	localClient := interactive.NewServiceAdapter(svc)
//	res := interactive.NewGreyScaleServiceClient(remoteClient, localClient, cfg.Threshold)
//	// 在这里监听
//	viper.OnConfigChange(func(in fsnotify.Event) {
//		cfg = Config{}
//		err1 := viper.UnmarshalKey("grpc.intr", cfg)
//		if err1 != nil {
//			l.Error("重新加载grpc.intr的配置失败", logger.Error(err1))
//			return
//		}
//		// 这边更新 Threshold
//		res.UpdateThreshold(cfg.Threshold)
//	})
//	return res
//}
