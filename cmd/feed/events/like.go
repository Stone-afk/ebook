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

const topicLikeEvent = "like_feed_event"

// LikeFeedEvent 由业务方定义，本服务做适配
type LikeFeedEvent struct {
	Uid   int64
	liked int
	biz   string
	bizId int64
}

type LikeFeedEventConsumer struct {
	client sarama.Client
	l      logger.Logger
	svc    service.FeedService
}

func NewLikeFeedEventConsumer(
	client sarama.Client,
	l logger.Logger,
	svc service.FeedService) *FollowerFeedEventConsumer {
	return &FollowerFeedEventConsumer{
		svc:    svc,
		client: client,
		l:      l,
	}
}

func (c *LikeFeedEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("like_feed_event",
		c.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{topicLikeEvent},
			saramax.NewHandler[LikeFeedEvent](c.l, c.Consume))
		if err != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (c *LikeFeedEventConsumer) Consume(msg *sarama.ConsumerMessage,
	evt LikeFeedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.svc.CreateFeedEvent(ctx, domain.FeedEvent{
		Type: service.FollowEventName,
		Ext: map[string]string{
			"uid":     strconv.FormatInt(evt.Uid, 10),
			"liked":   strconv.Itoa(evt.liked),
			"biz_idd": strconv.FormatInt(evt.bizId, 10),
			"biz":     evt.biz,
		},
	})
}
