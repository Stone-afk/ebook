package events

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/saramax"
	"ebook/cmd/reward/domain"
	"ebook/cmd/reward/service"
	"github.com/IBM/sarama"
	"strings"
	"time"
)

type PaymentEvent struct {
	BizTradeNO string
	Status     uint8
}

func (p PaymentEvent) ToDomainStatus() domain.RewardStatus {
	// 	PaymentStatusInit
	//	PaymentStatusSuccess
	//	PaymentStatusFailed
	//	PaymentStatusRefund
	switch p.Status {
	// 这里不能引用 payment 里面的定义，只能手写
	case 1:
		return domain.RewardStatusInit
	case 2:
		return domain.RewardStatusPayed
	case 3, 4:
		return domain.RewardStatusFailed
	default:
		return domain.RewardStatusUnknown
	}
}

type PaymentEventConsumer struct {
	client sarama.Client
	l      logger.Logger
	svc    service.RewardService
}

func NewPaymentEventConsumer(
	client sarama.Client,
	svc service.RewardService,
	l logger.Logger) *PaymentEventConsumer {
	return &PaymentEventConsumer{
		client: client,
		svc:    svc,
		l:      l,
	}
}

// Start 这边就是自己启动 goroutine 了
func (r *PaymentEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("reward",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{"payment_events"},
			saramax.NewHandler[PaymentEvent](r.l, r.Consume))
		if er != nil {
			r.l.Error("退出了消费循环异常", logger.Error(er))
		}
	}()
	return err
}

func (r *PaymentEventConsumer) Consume(
	msg *sarama.ConsumerMessage,
	evt PaymentEvent) error {
	// 不是我们的
	if !strings.HasPrefix(evt.BizTradeNO, "reward") {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return r.svc.UpdateReward(ctx, evt.BizTradeNO, evt.ToDomainStatus())
}
