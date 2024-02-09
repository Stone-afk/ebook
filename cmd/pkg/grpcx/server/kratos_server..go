package server

import (
	"ebook/cmd/pkg/grpcx"
	etcd "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"strconv"
)

type KratosServer struct {
	grpcx.Register
	Port       int
	EtcdAddrs  []string
	Name       string
	etcdClient *etcdv3.Client
	app        *kratos.App
}

func (s *KratosServer) Serve() error {
	grpcSrv := grpc.NewServer(
		grpc.Address(":"+strconv.Itoa(s.Port)),
		grpc.Middleware(recovery.Recovery()),
	)
	s.Registry(grpcSrv.Server)
	cli, err := etcdv3.New(etcdv3.Config{
		Endpoints: s.EtcdAddrs,
	})

	if err != nil {
		return err
	}
	s.etcdClient = cli
	r := etcd.New(s.etcdClient)
	s.app = kratos.New(
		kratos.Name(s.Name),
		kratos.Server(grpcSrv),
		kratos.Registrar(r))
	return s.app.Run()
}

// Close 你可以叫做 Shutdown
func (s *KratosServer) Close() error {
	if s.etcdClient != nil {
		err := s.etcdClient.Close()
		if err != nil {
			return err
		}
	}
	return s.app.Stop()
}
