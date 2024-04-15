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
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`

	Birthday      string `json:"birthday"`
	AboutMe       string `json:"about_me"`
	WechatOpenId  string `json:"wechat_open_id"`
	WechatUnionId string `json:"wechat_union_id"`

	// 创建时间
	Ctime int64 `json:"ctime"`
	// 更新时间
	Utime int64 `json:"utime"`
}

type IMUserEvent UserEvent

type SyncSearchEventProducer interface {
	ProduceStandardSyncEvent(ctx context.Context, evt UserEvent) error
	ProduceSyncEvent(ctx context.Context, evt UserEvent) error
}

type SyncIMEventProducer interface {
	ProduceIMEvent(ctx context.Context, evt IMUserEvent) error
}
