package grpc

import (
	"context"
	searchv1 "ebook/cmd/api/proto/gen/search/v1"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/service"
	"google.golang.org/grpc"
)

type SyncServiceServer struct {
	searchv1.UnimplementedSyncServiceServer
	syncSvc service.SyncService
}

func (s *SyncServiceServer) InputBizTags(ctx context.Context, request *searchv1.InputBizTagsRequest) (*searchv1.InputBizTagsResponse, error) {
	err := s.syncSvc.InputBizTags(ctx, toDomainBizTags(request.GetBiztags()))
	return &searchv1.InputBizTagsResponse{}, err
}

// InputUser 业务专属接口，你可以高度定制化
func (s *SyncServiceServer) InputUser(ctx context.Context, request *searchv1.InputUserRequest) (*searchv1.InputUserResponse, error) {
	err := s.syncSvc.InputUser(ctx, toDomainUser(request.GetUser()))
	return &searchv1.InputUserResponse{}, err
}

func (s *SyncServiceServer) InputArticle(ctx context.Context, request *searchv1.InputArticleRequest) (*searchv1.InputArticleResponse, error) {
	err := s.syncSvc.InputArticle(ctx, toDomainArticle(request.GetArticle()))
	return &searchv1.InputArticleResponse{}, err
}

func (s *SyncServiceServer) InputAny(ctx context.Context, req *searchv1.InputAnyRequest) (*searchv1.InputAnyResponse, error) {
	err := s.syncSvc.InputAny(ctx, req.IndexName, req.DocId, req.Data)
	return &searchv1.InputAnyResponse{}, err
}

func NewSyncServiceServer(syncSvc service.SyncService) *SyncServiceServer {
	return &SyncServiceServer{
		syncSvc: syncSvc,
	}
}

func (s *SyncServiceServer) Register(server grpc.ServiceRegistrar) {
	searchv1.RegisterSyncServiceServer(server, s)
}

func toDomainUser(vuser *searchv1.User) domain.User {
	return domain.User{
		Id:       vuser.Id,
		Email:    vuser.Email,
		Nickname: vuser.Nickname,
	}
}

func toDomainArticle(art *searchv1.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  art.Status,
		Content: art.Content,
		Tags:    art.Tags,
	}
}

func toDomainBizTags(bizTags *searchv1.BizTags) domain.BizTags {
	return domain.BizTags{
		Uid:   bizTags.Uid,
		Biz:   bizTags.Biz,
		BizId: bizTags.Bizid,
		Tags:  bizTags.Tags,
	}
}
