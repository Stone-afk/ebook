package events

import "context"

type SyncDataEvent struct {
	IndexName string
	DocID     string
	Data      string
}

type ArticleEvent struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Status  int32  `json:"status"`
	Content string `json:"content"`
}

type ReadEvent struct {
	Uid int64
	Aid int64
}

type ReadEventV1 struct {
	Uids []int64
	Aids []int64
}

type ReadEventProducer interface {
	ProduceReadEvent(ctx context.Context, evt ReadEvent) error
	//ProduceReadEventV1(ctx context.Context, v1 ReadEventV1) error
}

type SyncSearchEventProducer interface {
	ProduceStandardSyncEvent(ctx context.Context, evt ArticleEvent) error
	ProduceSyncEvent(ctx context.Context, evt ArticleEvent) error
	//ProduceReadEventV1(ctx context.Context, v1 ReadEventV1) error
}
