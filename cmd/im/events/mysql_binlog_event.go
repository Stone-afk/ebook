package events

import (
	"context"
	"ebook/cmd/im/domain"
	"ebook/cmd/im/service"
	"ebook/cmd/pkg/canalx"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/saramax"
	"github.com/IBM/sarama"
	"strconv"
	"time"
)

type SyncUserConsumer struct {
	client sarama.Client
	l      logger.Logger
	svc    service.UserService
}

func NewSyncUserConsumer(client sarama.Client,
	l logger.Logger, svc service.UserService) *SyncUserConsumer {
	return &SyncUserConsumer{
		client: client,
		l:      l,
		svc:    svc,
	}
}

func (r *SyncUserConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("open_im",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"open_im_binlog"},
			saramax.NewHandler[canalx.Message[User]](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *SyncUserConsumer) Consume(msg *sarama.ConsumerMessage,
	val canalx.Message[User]) error {
	// 因为共用了一个 topic，所以会有很多表的数据，不是自己的就不用管了
	// 只处理用户表的
	if val.Table != "users" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	for _, data := range val.Data {
		err := r.svc.Sync(ctx, domain.User{
			Nickname: data.Nickname,
			UserID:   strconv.FormatInt(data.Id, 10),
		})
		if err != nil {
			// 记录日志。
			continue
		}
	}
	return nil
}
