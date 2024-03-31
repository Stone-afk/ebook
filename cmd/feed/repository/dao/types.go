package dao

import "context"

// FeedPullEventDAO 拉模型
type FeedPullEventDAO interface {
	CreatePullEvent(ctx context.Context, event FeedPullEvent) error
	FindPullEvents(ctx context.Context, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error)
	FindPullEventsWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error)
}

type FeedPushEventDAO interface {
	// CreatePushEvents 创建推送事件
	CreatePushEvents(ctx context.Context, events []FeedPushEvent) error
	FindPushEvents(ctx context.Context, uid int64, timestamp, limit int64) ([]FeedPushEvent, error)
	FindPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp, limit int64) ([]FeedPushEvent, error)
}
