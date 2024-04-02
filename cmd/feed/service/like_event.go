package service

import (
	"context"
	"ebook/cmd/feed/domain"
	"ebook/cmd/feed/repository"
	"time"
)

const (
	LikeEventName = "like_event"
)

type LikeEventHandler struct {
	repo repository.FeedEventRepo
}

// CreateFeedEvent 中的 ext 里面至少需要三个 id
// liked int64: 被点赞的人
// liker int64：点赞的人
// bizId int64: 被点赞的东西
// biz: string
func (h *LikeEventHandler) CreateFeedEvent(ctx context.Context, ext domain.ExtendFields) error {
	liked, err := ext.Get("liked").AsInt64()
	if err != nil {
		return err
	}
	// 你可以考虑校验其它数据
	// 如果你用的是扩展表设计，那么这里就会调用自己业务的扩展表来存储数据
	// 如果你希望冗余存储数据，但是业务方又不愿意存，
	// 那么你在这里可以考虑回查业务获得一些数据
	return h.repo.CreatePushEvents(ctx, []domain.FeedEvent{{
		Uid:   liked,
		Type:  LikeEventName,
		Ctime: time.Now(),
		Ext:   ext,
	}})
}

func (h *LikeEventHandler) FindFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	return h.repo.FindPushEventsWithTyp(ctx, LikeEventName, uid, timestamp, limit)
}

func NewLikeEventHandler(repo repository.FeedEventRepo) Handler {
	return &LikeEventHandler{
		repo: repo,
	}
}
