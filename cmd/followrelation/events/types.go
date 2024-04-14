package events

import (
	"context"
	"ebook/cmd/pkg/typex"
)

type ExtendFields = typex.ExtendFields

// FollowerFeedEvent 由业务方定义，本服务做适配
type FollowerFeedEvent struct {
	Follower int64
	Followee int64
	Biz      string
	BizId    int64
}

type FeedEvent struct {
	Type     string
	Metadata ExtendFields
}

type FeedEventProducer interface {
	ProduceFeedEvent(ctx context.Context, evt FeedEvent) error
	ProduceStandardFeedEvent(ctx context.Context, evt FeedEvent) error
}
