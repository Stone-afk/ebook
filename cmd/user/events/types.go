package events

import "context"

type SyncDataEvent struct {
	IndexName string
	DocID     string
	Data      string
}

type UserEvent struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
}

type SyncSearchEventProducer interface {
	ProduceStandardSyncEvent(ctx context.Context, evt UserEvent) error
	ProduceSyncEvent(ctx context.Context, evt UserEvent) error
}
