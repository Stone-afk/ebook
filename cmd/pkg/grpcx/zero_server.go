package grpcx

import (
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"strconv"
)

type ZeroServer struct {
	Register
	Port      int
	EtcdAddrs []string
	Name      string
	server    *zrpc.RpcServer
}

func (s *ZeroServer) Serve() error {
	c := zrpc.RpcServerConf{
		// 这个是服务启动的地址
		ListenOn: ":" + strconv.Itoa(s.Port),
		Etcd: discov.EtcdConf{
			Hosts: s.EtcdAddrs,
			// 你的服务名
			Key: s.Name,
		},
	}
	server, err := zrpc.NewServer(c, func(server *grpc.Server) {
		// 吧你的业务注册到 server 里面
		s.Registry(grpc.NewServer())
	})

	if err != nil {
		return err
	}
	s.server = server
	s.server.Start()
	return nil
}
