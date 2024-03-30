package article

import (
	"context"
	"ebook/cmd/article/repository"
	"ebook/cmd/pkg/canalx"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/migrator"
	"ebook/cmd/pkg/migrator/events"
	"ebook/cmd/pkg/migrator/validator"
	"ebook/cmd/pkg/saramax"
	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"sync/atomic"
	"time"
)

const binlogTopic = "interactive_binlog"

type MySQLBinlogConsumer[T migrator.Entity] struct {
	client   sarama.Client
	l        logger.Logger
	table    string
	repo     *repository.CachedArticleRepository
	srcToDst *validator.CanalIncrValidator[T]
	dstToSrc *validator.CanalIncrValidator[T]
	dstFirst *atomic.Bool
}

func NewMySQLBinlogConsumer[T migrator.Entity](
	client sarama.Client,
	l logger.Logger,
	table string,
	src *gorm.DB,
	dst *gorm.DB,
	p events.Producer,
	repo *repository.CachedArticleRepository) *MySQLBinlogConsumer[T] {
	srcToDst := validator.NewCanalIncrValidator[T](src, dst, "SRC", l, p)
	dstToSrc := validator.NewCanalIncrValidator[T](src, dst, "DST", l, p)
	return &MySQLBinlogConsumer[T]{
		client: client, l: l,
		dstFirst: &atomic.Bool{},
		srcToDst: srcToDst,
		dstToSrc: dstToSrc,
		table:    table,
		repo:     repo}
}

func (r *MySQLBinlogConsumer[T]) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("migrator_incr",
		r.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{binlogTopic},
			saramax.NewHandler[canalx.Message[T]](r.l, r.Consume))
		if err != nil {
			r.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}

func (r *MySQLBinlogConsumer[T]) Consume(msg *sarama.ConsumerMessage,
	val canalx.Message[T]) error {
	dstFirst := r.dstFirst.Load()
	var v *validator.CanalIncrValidator[T]
	// db:
	//  src:
	//    dsn: "root:root@tcp(localhost:13316)/ebook"
	//  dst:
	//    dsn: "root:root@tcp(localhost:13316)/ebook_interactive"
	if dstFirst && val.Database == "ebook_interactive" {
		// 校验，用 dst 的来校验
		v = r.dstToSrc
	} else if !dstFirst && val.Database == "ebook" {
		v = r.srcToDst
	}
	if v != nil {
		for _, data := range val.Data {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := v.Validate(ctx, data.ID())
			cancel()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
