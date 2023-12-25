package article

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/saramax"
	"github.com/IBM/sarama"
	"time"
)

type HistoryReadEventConsumer struct {
	client sarama.Client
	repo   repository.HistoryRecordRepository
	l      logger.Logger
}

func NewHistoryReadEventConsumer(
	client sarama.Client,
	l logger.Logger,
	repo repository.HistoryRecordRepository) *HistoryReadEventConsumer {
	return &HistoryReadEventConsumer{
		client: client,
		l:      l,
		repo:   repo,
	}
}

func (r *HistoryReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", r.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{"read_article"},
			saramax.NewHandler[ReadEvent](r.l, r.Consume))
		if er != nil {
			r.l.Error("退出了消费循环异常", logger.Error(er))
		}
	}()
	return err
}

// Consume 这个不是幂等的
func (r *HistoryReadEventConsumer) Consume(msg *sarama.ConsumerMessage, evt ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return r.repo.AddRecord(ctx, domain.HistoryRecord{
		Uid:   evt.Uid,
		Biz:   "article",
		BizId: evt.Aid,
	})
}
