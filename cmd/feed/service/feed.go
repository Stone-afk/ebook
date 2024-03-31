package service

import (
	"context"
	"ebook/cmd/feed/domain"
	"ebook/cmd/feed/repository"
)

type feedService struct {
	repo       repository.FeedEventRepo
	handlerMap map[string]Handler
	//followClient followv1.FollowServiceClient
}

func (svc *feedService) CreateFeedEvent(ctx context.Context, feed domain.FeedEvent) error {
	//TODO implement me
	panic("implement me")
}

func (svc *feedService) GetFeedEventList(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	//TODO implement me
	panic("implement me")
}

func NewFeedService(repo repository.FeedEventRepo,
	//client followv1.FollowServiceClient,
	handlerMap map[string]Handler) FeedService {
	return &feedService{
		repo: repo,
		// 你可以注入那个 health client
		//followClient: client,
		handlerMap: handlerMap,
	}
}
