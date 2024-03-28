package validator

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/migrator/events"
	"gorm.io/gorm"
	"time"
)

type baseValidator struct {
	// 校验，以 XXX 为准，
	base *gorm.DB
	// 校验的是谁的数据
	target    *gorm.DB
	l         logger.Logger
	producer  events.Producer
	direction string
}

// 上报不一致的数据
func (v *baseValidator) notify(id int64, typ string) {
	// 这里我们要单独控制超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	evt := events.InconsistentEvent{
		Direction: v.direction,
		ID:        id,
		Type:      typ,
	}

	err := v.producer.ProduceInconsistentEvent(ctx, evt)
	if err != nil {
		v.l.Error("发送消息失败", logger.Error(err),
			logger.Field{Key: "event", Value: evt})
	}
}
