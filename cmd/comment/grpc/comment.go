package grpc

import (
	"context"
	commentv1 "ebook/cmd/api/proto/gen/comment/v1"
	"ebook/cmd/comment/domain"
	"ebook/cmd/comment/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
)

type CommentServiceServer struct {
	// 正常我都会组合这个
	commentv1.UnimplementedCommentServiceServer

	svc service.CommentService
}

func (s *CommentServiceServer) Register(server grpc.ServiceRegistrar) {
	commentv1.RegisterCommentServiceServer(server, s)
}

func (s *CommentServiceServer) GetMoreReplies(ctx context.Context, req *commentv1.GetMoreRepliesRequest) (*commentv1.GetMoreRepliesResponse, error) {
	cs, err := s.svc.GetMoreReplies(ctx, req.Rid, req.MaxId, req.Limit)
	if err != nil {
		return nil, err
	}
	return &commentv1.GetMoreRepliesResponse{
		Replies: s.toDTO(cs),
	}, nil
}

func (s *CommentServiceServer) GetCommentList(ctx context.Context, request *commentv1.CommentListRequest) (*commentv1.CommentListResponse, error) {
	minID := request.MinId
	// 第一次查询
	if minID <= 0 {
		minID = math.MaxInt64
	}
	domainComments, err := s.svc.
		GetCommentList(ctx,
			request.GetBiz(),
			request.GetBizid(),
			request.GetMinId(),
			request.GetLimit())
	if err != nil {
		return nil, err
	}
	return &commentv1.CommentListResponse{
		Comments: s.toDTO(domainComments),
	}, nil
}

func (s *CommentServiceServer) DeleteComment(ctx context.Context, request *commentv1.DeleteCommentRequest) (*commentv1.DeleteCommentResponse, error) {
	err := s.svc.DeleteComment(ctx, request.Id)
	return &commentv1.DeleteCommentResponse{}, err
}

func (s *CommentServiceServer) CreateComment(ctx context.Context, request *commentv1.CreateCommentRequest) (*commentv1.CreateCommentResponse, error) {
	err := s.svc.CreateComment(ctx, convertToDomain(request.GetComment()))
	return &commentv1.CreateCommentResponse{}, err
}

func (s *CommentServiceServer) toDTO(domainComments []domain.Comment) []*commentv1.Comment {
	rpcComments := make([]*commentv1.Comment, 0, len(domainComments))
	for _, domainComment := range domainComments {
		rpcComment := &commentv1.Comment{
			Id:      domainComment.Id,
			Uid:     domainComment.Commentator.ID,
			Biz:     domainComment.Biz,
			Bizid:   domainComment.BizID,
			Content: domainComment.Content,
			Ctime:   timestamppb.New(domainComment.CTime),
			Utime:   timestamppb.New(domainComment.UTime),
		}
		if domainComment.RootComment != nil {
			rpcComment.RootComment = &commentv1.Comment{
				Id: domainComment.RootComment.Id,
			}
		}
		if domainComment.ParentComment != nil {
			rpcComment.ParentComment = &commentv1.Comment{
				Id: domainComment.ParentComment.Id,
			}
		}
		rpcComments = append(rpcComments, rpcComment)
	}
	rpcCommentMap := make(map[int64]*commentv1.Comment, len(rpcComments))
	for _, rpcComment := range rpcComments {
		rpcCommentMap[rpcComment.Id] = rpcComment
	}
	for _, domainComment := range domainComments {
		rpcComment := rpcCommentMap[domainComment.Id]
		if domainComment.RootComment != nil {
			val, ok := rpcCommentMap[domainComment.RootComment.Id]
			if ok {
				rpcComment.RootComment = val
			}
		}
		if domainComment.ParentComment != nil {
			val, ok := rpcCommentMap[domainComment.ParentComment.Id]
			if ok {
				rpcComment.ParentComment = val
			}
		}
	}
	return rpcComments
}

func convertToDomain(comment *commentv1.Comment) domain.Comment {
	domainComment := domain.Comment{
		Id:      comment.Id,
		Biz:     comment.Biz,
		BizID:   comment.Bizid,
		Content: comment.Content,
		Commentator: domain.User{
			ID: comment.Uid,
		},
	}
	if comment.GetParentComment() != nil {
		domainComment.ParentComment = &domain.Comment{
			Id: comment.GetParentComment().GetId(),
		}
	}
	if comment.GetRootComment() != nil {
		domainComment.RootComment = &domain.Comment{
			Id: comment.GetRootComment().GetId(),
		}
	}
	return domainComment
}

func NewCommentServer(svc service.CommentService) *CommentServiceServer {
	return &CommentServiceServer{
		svc: svc,
	}
}
