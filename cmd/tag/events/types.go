package events

import "context"

type SyncDataEvent struct {
	IndexName string
	DocID     string
	Data      string
}

type BizTags struct {
	Uid   int64  `json:"uid"`
	Biz   string `json:"biz"`
	BizId int64  `json:"biz_id"`
	// 只传递 string
	Tags []string `json:"tags"`
}

type Producer interface {
	ProduceSyncEvent(ctx context.Context, data BizTags) error
	ProduceStandardSyncEvent(ctx context.Context, data BizTags) error
}
