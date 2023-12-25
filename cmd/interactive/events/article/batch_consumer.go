package article

import (
	"context"
	"ebook/cmd/interactive/repository"
	"ebook/cmd/internal/events/article"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/saramax"
	"github.com/IBM/sarama"
	"time"
)

type InteractiveReadEventBatchConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.Logger
}

func NewInteractiveReadEventBatchConsumer(client sarama.Client, repo repository.InteractiveRepository,
	l logger.Logger) *InteractiveReadEventBatchConsumer {
	return &InteractiveReadEventBatchConsumer{
		client: client,
		repo:   repo,
		l:      l,
	}
}

func (r *InteractiveReadEventBatchConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", r.client)
	if err != nil {
		return err
	}
	go func() {
		err = cg.Consume(context.Background(),
			[]string{"read_article"}, saramax.NewBatchHandler[article.ReadEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("线程退出消费, 循环异常", logger.Error(err))
		}
	}()
	return err
}

// Consume 这个不是幂等的
func (r *InteractiveReadEventBatchConsumer) Consume(msg []*sarama.ConsumerMessage,
	ts []article.ReadEvent) error {
	ids := make([]int64, 0, len(ts))
	bizs := make([]string, 0, len(ts))
	for _, evt := range ts {
		ids = append(ids, evt.Aid)
		bizs = append(bizs, "article")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := r.repo.BatchIncrReadCnt(ctx, bizs, ids)
	if err != nil {
		r.l.Error("批量增加阅读计数失败",
			logger.Field{Key: "ids", Value: ids},
			logger.Error(err))
	}
	return nil
}
