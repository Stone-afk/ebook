package service

import (
	"context"
	"ebook/cmd/feed/domain"
	"ebook/cmd/feed/repository"
	"time"
)

const (
	FollowEventName = "follow_event"
)

type FollowEventHandler struct {
	repo repository.FeedEventRepo
}

// CreateFeedEvent 创建跟随方式
// 如果 A 关注了 B，那么
// follower 就是 A
// followee 就是 B
func (h *FollowEventHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {
	followee, err := ext.Get("followee").AsInt64()
	if err != nil {
		return err
	}
	return h.repo.CreatePushEvents(ctx, []domain.FeedEvent{{
		Uid:   followee,
		Type:  FollowEventName,
		Ctime: time.Now(),
		Ext:   ext,
	}})
}

func (h *FollowEventHandler) FindFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	return h.repo.FindPushEventsWithTyp(ctx, FollowEventName, uid, timestamp, limit)
}

func NewFollowEventHandler(repo repository.FeedEventRepo) Handler {
	return &FollowEventHandler{
		repo: repo,
	}
}
