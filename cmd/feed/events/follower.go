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

const topicFollowerEvent = "follower_feed_event"

// FollowerFeedEvent 由业务方定义，本服务做适配
type FollowerFeedEvent struct {
	follower int64
	followee int64
	biz      string
	bizId    int64
}

type FollowerFeedEventConsumer struct {
	client sarama.Client
	l      logger.Logger
	svc    service.FeedService
}

func NewFollowerFeedEventConsumer(
	client sarama.Client,
	l logger.Logger,
	svc service.FeedService) *FollowerFeedEventConsumer {
	return &FollowerFeedEventConsumer{
		svc:    svc,
		client: client,
		l:      l,
	}
}

func (c *FollowerFeedEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("follower_feed_event",
		c.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicFollowerEvent},
			saramax.NewHandler[FollowerFeedEvent](c.l, c.Consume))
		if err != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (c *FollowerFeedEventConsumer) Consume(msg *sarama.ConsumerMessage,
	evt FollowerFeedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.svc.CreateFeedEvent(ctx, domain.FeedEvent{
		Type: service.FollowEventName,
		Ext: map[string]string{
			"followee": strconv.FormatInt(evt.followee, 10),
			"follower": strconv.FormatInt(evt.follower, 10),
			"biz_id":   strconv.FormatInt(evt.bizId, 10),
			"biz":      evt.biz,
		},
	})
}
