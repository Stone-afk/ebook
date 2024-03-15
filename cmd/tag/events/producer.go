package events

import "context"

type BizTags struct {
	Uid   int64  `json:"uid"`
	Biz   string `json:"biz"`
	BizId int64  `json:"biz_id"`
	// 只传递 string
	Tags []string `json:"tags"`
}

type Producer interface {
	ProduceSyncEvent(ctx context.Context, data BizTags) error
}
