package events

import (
	"context"
	"ebook/cmd/feed/domain"
	"ebook/cmd/feed/service"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/saramax"
	"github.com/IBM/sarama"
	"strconv"
	"time"
)

const topicArticleEvent = "article_feed_event"

// ArticleFeedEvent 由业务方定义，本服务做适配
type ArticleFeedEvent struct {
	uid int64
	aid int64
}

type ArticleEventConsumer struct {
	client sarama.Client
	l      logger.Logger
	svc    service.FeedService
}

func NewArticleEventConsumer(
	client sarama.Client,
	l logger.Logger,
	svc service.FeedService) *ArticleEventConsumer {
	return &ArticleEventConsumer{
		svc:    svc,
		client: client,
		l:      l,
	}
}

// Start 这边就是自己启动 goroutine 了
func (c *ArticleEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("article_feed_event",
		c.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicArticleEvent},
			saramax.NewHandler[ArticleFeedEvent](c.l, c.Consume))
		if err != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (c *ArticleEventConsumer) Consume(msg *sarama.ConsumerMessage,
	evt ArticleFeedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.svc.CreateFeedEvent(ctx, domain.FeedEvent{
		Type: service.ArticleEventName,
		Ext: map[string]string{
			"uid": strconv.FormatInt(evt.uid, 10),
			"aid": strconv.FormatInt(evt.uid, 10),
		},
	})
}
