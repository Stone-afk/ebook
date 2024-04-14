package like

import (
	"context"
	"ebook/cmd/pkg/typex"
)

// LikeFeedEvent 由业务方定义，本服务做适配
type LikeFeedEvent struct {
	Uid   int64
	Liked int
	Biz   string
	BizId int64
}

type FeedEvent struct {
	Type     string
	Metadata ExtendFields
}

type ExtendFields = typex.ExtendFields

type FeedEventProducer interface {
	ProduceFeedEvent(ctx context.Context, evt FeedEvent) error
	ProduceStandardFeedEvent(ctx context.Context, evt FeedEvent) error
}
