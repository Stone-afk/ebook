package service

import (
	"context"
	"ebook/cmd/feed/domain"
	"ebook/cmd/feed/repository"
)

const (
	LikeEventName = "like_event"
)

type LikeEventHandler struct {
	repo repository.FeedEventRepo
}

func (h *LikeEventHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {
	//TODO implement me
	panic("implement me")
}

func (h *LikeEventHandler) FindFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	//TODO implement me
	panic("implement me")
}

func NewLikeEventHandler(repo repository.FeedEventRepo) Handler {
	return &LikeEventHandler{
		repo: repo,
	}
}
