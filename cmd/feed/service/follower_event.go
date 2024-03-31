package service

import (
	"context"
	"ebook/cmd/feed/domain"
	"ebook/cmd/feed/repository"
)

const (
	FollowEventName = "follow_event"
)

type FollowEventHandler struct {
	repo repository.FeedEventRepo
}

func (h *FollowEventHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {
	//TODO implement me
	panic("implement me")
}

func (h *FollowEventHandler) FindFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	//TODO implement me
	panic("implement me")
}

func NewFollowEventHandler(repo repository.FeedEventRepo) Handler {
	return &FollowEventHandler{
		repo: repo,
	}
}
