package grpcx

import (
	"ebook/cmd/pkg/logger"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type Server struct {
	*grpc.Server
	//Addr string
	L      logger.Logger
	em     endpoints.Manager
	client *etcdv3.Client

	Port      int
	EtcdAddrs []string
	Name      string
	kaCancel  func()
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		return err
	}
	// 这边会阻塞，类似与 gin.Run
	return s.Server.Serve(l)
}

func (s *Server) register() error {
	panic("")
}

// Close 你可以叫做 Shutdown
func (s *Server) Close() error {
	panic("")
}
