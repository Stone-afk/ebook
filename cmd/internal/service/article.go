package service

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository"
	"ebook/cmd/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	PublishV1(ctx context.Context, art domain.Article) (int64, error)
}

type articleService struct {
	// 1. 在 service 这一层使用两个 repository
	authorRepo repository.ArticleAuthorRepository
	readerRepo repository.ArticleReaderRepository

	// 2. 在 repo 里面处理制作库和线上库
	// 1 和 2 是互斥的，不会同时存在
	repo repository.ArticleRepository
	log  logger.Logger
}

func NewArticleService(repo repository.ArticleRepository, l logger.Logger) ArticleService {
	return &articleService{
		repo: repo,
		log:  l,
	}
}

func NewArticleServiceV1(
	authorRepo repository.ArticleAuthorRepository,
	readerRepo repository.ArticleReaderRepository,
	l logger.Logger) ArticleService {
	return &articleService{
		authorRepo: authorRepo,
		readerRepo: readerRepo,
		log:        l,
	}
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
