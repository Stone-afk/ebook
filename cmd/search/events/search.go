package events

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/saramax"
	"ebook/cmd/search/service"
	"github.com/IBM/sarama"
	"time"
)

const topicSyncData = "sync_data_event"

type SyncDataEvent struct {
	IndexName string
	DocID     string
	Data      string
	// 假如说用于同步 user
	// IndexName = user_index
	// DocID = "123"
	// Data = {"id": 123, "email":xx, nickname: ""}
}

type SyncDataEventConsumer struct {
	svc    service.SyncService
	client sarama.Client
	l      logger.Logger
}

func (c *SyncDataEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("search_sync_data",
		c.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicSyncData},
			saramax.NewHandler[SyncDataEvent](c.l, c.Consume))
		if err != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (c *SyncDataEventConsumer) Consume(sg *sarama.ConsumerMessage,
	evt SyncDataEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 在这里执行转化
	return c.svc.InputAny(ctx, evt.IndexName, evt.DocID, evt.Data)
}

func NewSyncDataEventConsumer(svc service.SyncService, client sarama.Client, l logger.Logger) saramax.Consumer {
	return &SyncDataEventConsumer{
		svc:    svc,
		client: client,
		l:      l,
	}
}
