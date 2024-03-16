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

const topicSyncUser = "sync_user_event"

type UserEvent struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
}

type UserConsumer struct {
	syncSvc service.SyncService
	client  sarama.Client
	l       logger.Logger
}

func NewUserConsumer(client sarama.Client,
	l logger.Logger,
	svc service.SyncService) *UserConsumer {
	return &UserConsumer{
		syncSvc: svc,
		client:  client,
		l:       l,
	}
}

func (c *UserConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("sync_user",
		c.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{topicSyncUser},
			saramax.NewHandler[UserEvent](c.l, c.Consume))
		if er != nil {
			c.l.Error("退出了消费循环异常", logger.Error(er))
		}
	}()
	return err
}

func (c *UserConsumer) Consume(sg *sarama.ConsumerMessage,
	evt UserEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.syncSvc.InputUser(ctx, c.toDomain(evt))
}

func (c *UserConsumer) toDomain(evt UserEvent) domain.User {
	return domain.User{
		Id:       evt.Id,
		Email:    evt.Email,
		Nickname: evt.Nickname,
		Phone:    evt.Phone,
	}
}
