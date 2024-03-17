package events

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/saramax"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/service"
	"github.com/IBM/sarama"
	"time"
)

const topicSyncTag = "sync_tag_event"

type TagEvent struct {
	Uid   int64  `json:"uid"`
	Biz   string `json:"biz"`
	BizId int64  `json:"biz_id"`
	// 只传递 string
	Tags []string `json:"tags"`
}

type TagConsumer struct {
	syncSvc service.SyncService
	client  sarama.Client
	l       logger.Logger
}

func NewTagConsumer(client sarama.Client,
	l logger.Logger,
	svc service.SyncService) *ArticleConsumer {
	return &ArticleConsumer{
		syncSvc: svc,
		client:  client,
		l:       l,
	}
}

func (c *TagConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("sync_tag",
		c.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{topicSyncTag},
			saramax.NewHandler[TagEvent](c.l, c.Consume))
		if er != nil {
			c.l.Error("退出了消费循环异常", logger.Error(er))
		}
	}()
	return err
}

func (c *TagConsumer) Consume(sg *sarama.ConsumerMessage,
	evt TagEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.syncSvc.InputBizTags(ctx, c.toDomain(evt))
}

func (c *TagConsumer) toDomain(event TagEvent) domain.BizTags {
	return domain.BizTags{
		Uid:   event.Uid,
		Biz:   event.Biz,
		BizId: event.BizId,
		// 只传递 string
		Tags: event.Tags,
	}
}
