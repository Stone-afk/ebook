package service

import (
	"context"
	followv1 "ebook/cmd/api/proto/gen/followrelation/v1"
	"ebook/cmd/feed/domain"
	"ebook/cmd/feed/repository"
)

const (
	ArticleEventName = "article_event"
	threshold        = 4
	//threshold        = 32
)

type ArticleEventHandler struct {
	repo         repository.FeedEventRepo
	followClient followv1.FollowServiceClient
}

func (h *ArticleEventHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {
	//TODO implement me
	panic("implement me")
}

func (h *ArticleEventHandler) FindFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	//TODO implement me
	panic("implement me")
}

func NewArticleEventHandler(repo repository.FeedEventRepo, client followv1.FollowServiceClient) Handler {
	return &ArticleEventHandler{
		repo:         repo,
		followClient: client,
	}
}
