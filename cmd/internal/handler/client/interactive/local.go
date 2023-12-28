package interactive

import (
	"context"
	intrv1 "ebook/cmd/api/proto/gen/intr/v1"
	domain2 "ebook/cmd/interactive/domain"
	"ebook/cmd/interactive/service"
	"google.golang.org/grpc"
)

// ServiceAdapter 将一个本地实现伪装成一个 gRPC 客户端
type ServiceAdapter struct {
	svc service.InteractiveService
}

func NewServiceAdapter(svc service.InteractiveService) *ServiceAdapter {
	return &ServiceAdapter{svc: svc}
}

func (s *ServiceAdapter) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	err := s.svc.IncrReadCnt(ctx, in.GetBiz(), in.GetBizId())
	return &intrv1.IncrReadCntResponse{}, err
}

func (s *ServiceAdapter) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	err := s.svc.Like(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.LikeResponse{}, err
}

func (s *ServiceAdapter) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	err := s.svc.CancelLike(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.CancelLikeResponse{}, err
}

func (s *ServiceAdapter) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	err := s.svc.Collect(ctx, in.GetBiz(), in.GetBizId(), in.GetUid(), in.GetCid())
	return &intrv1.CollectResponse{}, err
}

func (s *ServiceAdapter) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	intr, err := s.svc.Get(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResponse{
		Intr: s.toDTO(intr),
	}, nil
}

func (s *ServiceAdapter) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	res, err := s.svc.GetByIds(ctx, in.GetBiz(), in.GetIds())
	if err != nil {
		return nil, err
	}
	m := make(map[int64]*intrv1.Interactive, len(res))
	for k, v := range res {
		m[k] = s.toDTO(v)
	}
	return &intrv1.GetByIdsResponse{
		Intrs: m,
	}, nil
}

// DTO data transfer object
func (s *ServiceAdapter) toDTO(intr domain2.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		CollectCnt: intr.CollectCnt,
		Collected:  intr.Collected,
		LikeCnt:    intr.LikeCnt,
		Liked:      intr.Liked,
		ReadCnt:    intr.ReadCnt,
	}
}
