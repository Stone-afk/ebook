package main

import (
	"context"
	intrv1 "ebook/cmd/api/proto/gen/intr/v1"
	"ebook/cmd/pkg/netx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"
)

type EtcdTestSuite struct {
	suite.Suite
	client *etcdv3.Client
}

func (s *EtcdTestSuite) SetupSuite() {
	client, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:12379"},
	})
	require.NoError(s.T(), err)
	s.client = client
}

func (s *EtcdTestSuite) TestClient() {
	// build 注册中心
	bd, err := resolver.NewBuilder(s.client)
	require.NoError(s.T(), err)
	// 连接注册中心
	cc, err := grpc.Dial("etcd:///service/interactive", grpc.WithResolvers(bd),
		//grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		//	ctx = context.WithValue(ctx, "req", req)
		//	return invoker(ctx, method, req, reply, cc)
		//}),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	client := intrv1.NewInteractiveServiceClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//ctx = context.WithValue(ctx, "balancer-key", 123)
	resp, err := client.IncrReadCnt(ctx, &intrv1.IncrReadCntRequest{
		Biz:   "article",
		BizId: 123,
	})
	require.NoError(s.T(), err)
	s.T().Log(resp.String())
	time.Sleep(time.Minute)
}

func (s *EtcdTestSuite) TestServer() {
	l, err := net.Listen("tcp", ":8090")
	require.NoError(s.T(), err)
	// endpoint 以服务为维度。一个服务一个 Manager
	em, err := endpoints.NewManager(s.client, "service/interactive")
	require.NoError(s.T(), err)
	addr := netx.GetOutboundIP() + ":8090"
	// key 是指这个实例的 key
	// 如果有 instance id，用 instance id，如果没有，本机 IP + 端口
	// 端口一般是从配置文件里面读
	key := "service/interactive/" + addr
	//... 在这一步之前完成所有的启动的准备工作，包括缓存预加载之类的事情

	// 这个 ctx 是控制创建租约的超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// ttl 是租期
	// 秒作为单位
	// 过了 1/3（还剩下 2/3 的时候）就续约
	var ttl int64 = 30
	leaseResp, err := s.client.Grant(ctx, ttl)
	require.NoError(s.T(), err)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
		Addr: addr,
		Metadata: map[string]any{
			"weight": 100,
			"cpu":    90,
		},
	}, etcdv3.WithLease(leaseResp.ID))
	require.NoError(s.T(), err)
	kaCtx, kaCancel := context.WithCancel(context.Background())

	// 心跳
	go func() {
		ch, err1 := s.client.KeepAlive(kaCtx, leaseResp.ID)
		require.NoError(s.T(), err1)
		for kaResp := range ch {
			// 正常就是打印一下 DEBUG 日志啥的
			s.T().Log(kaResp.String(), time.Now().String())
		}
	}()

	// 续约
	go func() {
		ticker := time.NewTicker(time.Second)
		for now := range ticker.C {
			ctx1, cancel1 := context.WithTimeout(context.Background(), time.Second)
			// AddEndpoint 是一个覆盖的语义。也就是说，如果你这边已经有这个 key 了，就覆盖
			// upsert，set
			err = em.AddEndpoint(ctx1, key, endpoints.Endpoint{
				Addr: addr,
				// 你的分组信息，权重信息，机房信息
				// 以及动态判定负载的时候，可以把你的负载信息也写到这里
				Metadata: map[string]any{
					"weight": 200,
					"time":   now.String(),
				},
			}, etcdv3.WithLease(leaseResp.ID))
			if err != nil {
				s.T().Log(err)
			}
			//// 最开始，以为 Update 是需要自己手工删了，然后再加上去
			//em.Update(ctx1, []*endpoints.UpdateWithOpts{
			//	{
			//		Update: endpoints.Update{
			//			// Op 只有 Add 和 Delete
			//			// 没有 Update
			//			Op:  endpoints.Delete,
			//			Key: key,
			//		},
			//	},
			//{
			//	Update: endpoints.Update{
			//		Op:  endpoints.Add,
			//		Key: key,
			//		Endpoint: endpoints.Endpoint{
			//			Addr: addr,
			//			// 你的分组信息，权重信息，机房信息
			//			// 以及动态判定负载的时候，可以把你的负载信息也写到这里
			//			Metadata: now.String(),
			//		},
			//	},
			//},
			//})
			cancel1()
		}
	}()
	server := grpc.NewServer()
	intrv1.RegisterInteractiveServiceServer(server, &Server{})
	err = server.Serve(l)
	s.T().Log(err)
	// 要退出了，正常退出
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 我要先取消续约
	kaCancel()
	// 退出阶段，先从注册中心里面删了自己
	err = em.DeleteEndpoint(ctx, key)
	// 关掉客户端
	_ = s.client.Close()
	server.GracefulStop()
}

func TestEtcd(t *testing.T) {
	suite.Run(t, new(EtcdTestSuite))
}

type Server struct {
	intrv1.UnimplementedInteractiveServiceServer
}

var _ intrv1.InteractiveServiceServer = &Server{}

func (s *Server) IncrReadCnt(ctx context.Context, request *intrv1.IncrReadCntRequest) (*intrv1.IncrReadCntResponse, error) {
	return &intrv1.IncrReadCntResponse{}, nil
}
