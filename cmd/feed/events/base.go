package events

import (
	"context"
	"ebook/cmd/feed/domain"
	"ebook/cmd/feed/service"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/saramax"
	"github.com/IBM/sarama"
	"time"
)

const topicFeedEvent = "feed_event"

type FeedEvent struct {
	Type     string
	Metadata map[string]string
}

type FeedEventConsumer struct {
	client sarama.Client
	l      logger.Logger
	svc    service.FeedService
}

func NewFeedEventConsumer(
	client sarama.Client,
	l logger.Logger,
	svc service.FeedService) *FeedEventConsumer {
	return &FeedEventConsumer{
		svc:    svc,
		client: client,
		l:      l,
	}
}

func (c *FeedEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("feed_event",
		c.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicFeedEvent},
			saramax.NewHandler[FeedEvent](c.l, c.Consume))
		if err != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (c *FeedEventConsumer) Consume(msg *sarama.ConsumerMessage,
	evt FeedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.svc.CreateFeedEvent(ctx, domain.FeedEvent{
		Type: evt.Type,
		Ext:  evt.Metadata,
	})
}
