package article

import (
	"ebook/cmd/article/repository"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/migrator"
	"github.com/IBM/sarama"
	"sync/atomic"
)

type MySQLBinlogConsumer[T migrator.Entity] struct {
	client   sarama.Client
	l        logger.Logger
	table    string
	repo     *repository.CachedArticleRepository
	dstFirst *atomic.Bool
}
