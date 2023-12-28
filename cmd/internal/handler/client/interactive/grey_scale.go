package interactive

import (
	"context"
	intrv1 "ebook/cmd/api/proto/gen/intr/v1"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"google.golang.org/grpc"
)

type GreyScaleServiceClient struct {
	remote intrv1.InteractiveServiceClient
	local  intrv1.InteractiveServiceClient
	// 怎么控制流量呢？
	// 如果一个请求过来，该怎么控制它去调用本地，还是调用远程呢？
	// 用随机数 + 阈值的小技巧
	threshold *atomicx.Value[int32]
}

func NewGreyScaleServiceClient(remote intrv1.InteractiveServiceClient, local intrv1.InteractiveServiceClient, threshold int32) *GreyScaleServiceClient {
	return &GreyScaleServiceClient{
		remote:    remote,
		local:     local,
		threshold: atomicx.NewValueOf(threshold),
	}
}

func (c *GreyScaleServiceClient) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	panic("")
}

func (c *GreyScaleServiceClient) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	panic("")
}

func (c *GreyScaleServiceClient) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	panic("")
}

func (c *GreyScaleServiceClient) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	panic("")
}

func (c *GreyScaleServiceClient) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	panic("")
}

func (c *GreyScaleServiceClient) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	panic("")
}

func (c *GreyScaleServiceClient) UpdateThreshold(threshold int32) {
	c.threshold.Store(threshold)
}
