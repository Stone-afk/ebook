package interactive

import (
	"context"
	intrv1 "ebook/cmd/api/proto/gen/intr/v1"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	bd, err := resolver.NewBuilder(s.client)
	require.NoError(s.T(), err)
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

}

func TestEtcd(t *testing.T) {
	suite.Run(t, new(EtcdTestSuite))
}
