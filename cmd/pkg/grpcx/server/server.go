package server

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/netx"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
)

type Server struct {
	*grpc.Server
	//Addr string
	L          logger.Logger
	em         endpoints.Manager
	etcdClient *etcdv3.Client

	key       string
	Port      int
	EtcdAddrs []string
	Name      string
	kaCancel  func()
	// ETCD 服务注册租约 TTL
	EtcdTTL int64
}

func NewGRPCXServer(grpcSvc *grpc.Server,
	etcdClient *etcdv3.Client,
	l logger.Logger,
	port int,
	serverName string,
	etcdTTL int64) *Server {
	return &Server{
		L:          l,
		Server:     grpcSvc,
		etcdClient: etcdClient,
		Port:       port,
		Name:       serverName,
		EtcdTTL:    etcdTTL,
	}
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		return err
	}
	err = s.register()
	if err != nil {
		return err
	}
	// 就是直接启动，现在要嵌入服务注册过程
	// 这边会阻塞，类似与 gin.Run
	return s.Server.Serve(l)
}

func (s *Server) register() error {
	if s.etcdClient == nil {
		client, err := etcdv3.New(etcdv3.Config{
			Endpoints: s.EtcdAddrs,
		})
		if err != nil {
			return err
		}
		s.etcdClient = client
	}
	// endpoint 以服务为维度。一个服务一个 Manager
	em, err := endpoints.NewManager(s.etcdClient, "service/"+s.Name)
	if err != nil {
		return err
	}

	addr := netx.GetOutboundIP() + ":" + strconv.Itoa(s.Port)
	key := "service/" + s.Name + "/" + addr
	s.key = key
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//// 你可以做成配置的
	//var ttl int64 = 30
	leaseResp, err := s.etcdClient.Grant(ctx, s.EtcdTTL)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
	}, etcdv3.WithLease(leaseResp.ID))

	kaCtx, kaCancel := context.WithCancel(context.Background())
	s.kaCancel = kaCancel
	ch, err := s.etcdClient.KeepAlive(kaCtx, leaseResp.ID)
	go func() {
		for kaResp := range ch {
			// 正常就是打印一下 DEBUG 日志啥的
			s.L.Debug(kaResp.String())
		}
	}()
	return nil
}

// Close 你可以叫做 Shutdown
func (s *Server) Close() error {
	if s.kaCancel != nil {
		s.kaCancel()
	}
	if s.em != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := s.em.DeleteEndpoint(ctx, s.key)
		if err != nil {
			return err
		}
	}
	if s.etcdClient != nil {
		err := s.etcdClient.Close()
		if err != nil {
			return err
		}
	}
	s.Server.GracefulStop()
	return nil
}
