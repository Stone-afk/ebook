package repository

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository/dao/article"
	"ebook/cmd/pkg/logger"
	"gorm.io/gorm"
)

// repository 还是要用来操作缓存和DAO
// 事务概念应该在 DAO 这一层

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	// Sync 本身要求先保存到制作库，再同步到线上库
	Sync(ctx context.Context, art domain.Article) (int64, error)
	// SyncStatus 仅仅同步状态
	SyncStatus(ctx context.Context, uid, id int64, status domain.ArticleStatus) error
}

type articleRepository struct {
	dao article.ArticleDAO

	// SyncV1 用
	authorDAO article.ArticleAuthorDAO
	readerDAO article.ArticleReaderDAO

	// SyncV2 用
	db *gorm.DB
	l  logger.Logger
}

func NewArticleRepository(dao article.ArticleDAO, l logger.Logger) ArticleRepository {
	return &articleRepository{
		dao: dao,
		l:   l,
	}
}

func NewArticleRepositoryV1(authorDAO article.ArticleAuthorDAO,
	readerDAO article.ArticleReaderDAO) ArticleRepository {
	return &articleRepository{
		authorDAO: authorDAO,
		readerDAO: readerDAO,
	}
}

func NewArticleRepositoryV2(db *gorm.DB, l logger.Logger) ArticleRepository {
	return &articleRepository{
		db: db,
		l:  l,
	}
}

func (repo *articleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	id, err := repo.dao.Sync(ctx, repo.toEntity(art))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := repo.dao.Sync(ctx, repo.toEntity(art))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *articleRepository) SyncStatus(ctx context.Context, uid, id int64, status domain.ArticleStatus) error {
	return repo.dao.SyncStatus(ctx, uid, id, status.ToUint8())
}

func (repo *articleRepository) Update(ctx context.Context, art domain.Article) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(art))
}

func (repo *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(art))
}

func (repo *articleRepository) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   uint8(art.Status),
	}
}

func (repo *articleRepository) toDomain(art article.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Status: domain.ArticleStatus(art.Status),
	}
}
