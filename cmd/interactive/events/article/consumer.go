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

type InteractiveReadEventConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.Logger
}

func NewInteractiveReadEventConsumer(
	client sarama.Client,
	repo repository.InteractiveRepository,
	l logger.Logger) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{
		client: client,
		l:      l,
		repo:   repo,
	}
}

func (r *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient(
		"interactive", r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(
			context.Background(),
			[]string{"read_article"}, saramax.NewHandler[article.ReadEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

// Consume 这个不是幂等的
func (r *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t article.ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return r.repo.IncrReadCnt(ctx, "article", t.Aid)
}
