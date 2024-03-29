package article

import (
	"context"
	"ebook/cmd/interactive/repository"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/saramax"
	"github.com/IBM/sarama"
	"time"
)

const topicReadEvent = "article_read_event"

var _ saramax.Consumer = &InteractiveReadEventConsumer{}

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
		err = cg.Consume(
			context.Background(),
			[]string{topicReadEvent}, saramax.NewHandler[ReadEvent](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *InteractiveReadEventConsumer) StartBatch() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", r.client)
	if err != nil {
		return err
	}
	go func() {
		err = cg.Consume(context.Background(),
			[]string{topicReadEvent}, saramax.NewBatchHandler[ReadEvent](r.l, r.BatchConsume))
		if err != nil {
			r.l.Error("线程退出消费, 循环异常", logger.Error(err))
		}
	}()
	return err
}

// Consume 这个不是幂等的
func (r *InteractiveReadEventConsumer) Consume(msg *sarama.ConsumerMessage, t ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return r.repo.IncrReadCnt(ctx, "article", t.Aid)
}

func (r *InteractiveReadEventConsumer) BatchConsume(msgs []*sarama.ConsumerMessage,
	evts []ReadEvent) error {
	ids := make([]int64, 0, len(msgs))
	bizs := make([]string, 0, len(msgs))
	for _, evt := range evts {
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
