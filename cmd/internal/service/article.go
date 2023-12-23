package service

import (
	"context"
	"ebook/cmd/internal/domain"
	events "ebook/cmd/internal/events/article"
	"ebook/cmd/internal/repository"
	"ebook/cmd/pkg/logger"
	"time"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/service/article.go -package=svcmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/service/mocks/article.mock.go
type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
	List(ctx context.Context, authorId int64, offset, limit int) ([]domain.Article, error)
	// ListPub 只会取 start 七天内的数据
	ListPub(ctx context.Context, uTime time.Time, offset, limit int) ([]domain.Article, error)
	GetPublishedById(ctx context.Context, id int64, userId int64) (domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
}

type articleService struct {
	// 1. 在 service 这一层使用两个 repository
	authorRepo repository.ArticleAuthorRepository
	readerRepo repository.ArticleReaderRepository

	// 2. 在 repo 里面处理制作库和线上库
	// 1 和 2 是互斥的，不会同时存在
	repo     repository.ArticleRepository
	producer events.Producer
	log      logger.Logger

	ch chan readInfo
}

func NewArticleService(repo repository.ArticleRepository,
	l logger.Logger, producer events.Producer) ArticleService {
	return &articleService{
		repo:     repo,
		log:      l,
		producer: producer,
	}
}

//func NewArticleServiceV2(repo repository.ArticleRepository,
//	l logger.Logger,
//	producer events.Producer) ArticleService {
//	ch := make(chan readInfo, 10)
//	go func() {
//		for {
//			uids := make([]int64, 0, 10)
//			aids := make([]int64, 0, 10)
//			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//			for i := 0; i < 10; i++ {
//				select {
//				case info, ok := <-ch:
//					if !ok {
//						cancel()
//						return
//					}
//					uids = append(uids, info.uid)
//					aids = append(aids, info.aid)
//				case <-ctx.Done():
//					break
//				}
//			}
//			cancel()
//			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
//			err := producer.ProduceReadEventV1(ctx, events.ReadEventV1{
//				Uids: uids,
//				Aids: aids,
//			})
//			if err == nil {
//				l.Error("发送读者阅读事件失败")
//			}
//			cancel()
//		}
//	}()
//	return &articleService{
//		repo:     repo,
//		producer: producer,
//		log:      l,
//		ch:       ch,
//	}
//}

type readInfo struct {
	uid int64
	aid int64
}

func NewArticleServiceV1(
	authorRepo repository.ArticleAuthorRepository,
	readerRepo repository.ArticleReaderRepository,
	producer events.Producer,
	l logger.Logger) ArticleService {
	return &articleService{
		authorRepo: authorRepo,
		readerRepo: readerRepo,
		log:        l,
		producer:   producer,
	}
}

func (svc *articleService) ListPub(ctx context.Context,
	uTime time.Time, offset, limit int) ([]domain.Article, error) {
	return svc.repo.ListPub(ctx, uTime, offset, limit)
}

func (svc *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return svc.repo.GetById(ctx, id)
}

func (svc *articleService) GetPublishedById(ctx context.Context, id, userId int64) (domain.Article, error) {
	art, err := svc.repo.GetPublishedById(ctx, id)
	if err == nil {
		go func() {
			// 生产者也可以通过改批量来提高性能
			er := svc.producer.ProduceReadEvent(
				ctx,
				events.ReadEvent{
					// 即便你的消费者要用 art 的里面的数据，
					// 让它去查询，你不要在 event 里面带
					Uid: userId,
					Aid: id,
				})
			if er == nil {
				svc.log.Error("发送读者阅读事件失败")
			}
		}()
	}
	return art, err
}

func (svc *articleService) List(ctx context.Context, authorId int64, offset, limit int) ([]domain.Article, error) {
	return svc.repo.List(ctx, authorId, offset, limit)
}

func (svc *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	// art.Status = domain.ArticleStatusPrivate 然后直接把整个 art 往下传
	return svc.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func (svc *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := svc.repo.Update(ctx, art)
		return art.Id, err
	}
	return svc.repo.Create(ctx, art)
}

func (svc *articleService) Publish(ctx context.Context,
	art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return svc.repo.Sync(ctx, art)
}

// PublishV1 基于使用两种 repository 的写法
func (svc *articleService) PublishV1(ctx context.Context,
	art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	// 这一段逻辑其实就是 Save
	if art.Id == 0 {
		id, err = svc.authorRepo.Create(ctx, art)
	} else {
		err = svc.authorRepo.Update(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	// 保持制作库和线上库的 ID 是一样的。
	art.Id = id
	for i := 0; i < 3; i++ {
		err = svc.readerRepo.Save(ctx, art)
		if err == nil {
			break
		}
		svc.log.Error("部分失败：保存数据到线上库失败",
			logger.Int64("art_id", id),
			logger.Error(err))
	}
	// 在接入了 metrics 或者 tracing 之后，
	// 这边要进一步记录必要的DEBUG信息。
	if err != nil {
		svc.log.Error("全部失败：保存数据到线上库重试都失败了",
			logger.Int64("art_id", id),
			logger.Error(err))
	}
	return id, nil
}
