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

const topicSyncArticle = "sync_article_event"

type ArticleEvent struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Status  int32  `json:"status"`
	Content string `json:"content"`
}

type ArticleConsumer struct {
	syncSvc service.SyncService
	client  sarama.Client
	l       logger.Logger
}

func NewArticleConsumer(client sarama.Client,
	l logger.Logger,
	svc service.SyncService) *ArticleConsumer {
	return &ArticleConsumer{
		syncSvc: svc,
		client:  client,
		l:       l,
	}
}

func (c *ArticleConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("sync_article",
		c.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{topicSyncArticle},
			saramax.NewHandler[ArticleEvent](c.l, c.Consume))
		if er != nil {
			c.l.Error("退出了消费循环异常", logger.Error(er))
		}
	}()
	return err
}

func (c *ArticleConsumer) Consume(sg *sarama.ConsumerMessage,
	evt ArticleEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.syncSvc.InputArticle(ctx, c.toDomain(evt))
}

func (c *ArticleConsumer) toDomain(article ArticleEvent) domain.Article {
	return domain.Article{
		Id:      article.Id,
		Title:   article.Title,
		Status:  article.Status,
		Content: article.Content,
	}
}
