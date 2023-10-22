package repository

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository/dao/article"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, art domain.Article) error
}

type articleReaderRepository struct {
	dao article.ArticleReaderDAO
}

func NewArticleReaderRepository(dao article.ArticleReaderDAO) ArticleReaderRepository {
	return &articleReaderRepository{
		dao: dao,
	}
}

func (repo *articleReaderRepository) Save(ctx context.Context, art domain.Article) error {
	return repo.dao.Upsert(ctx, repo.toEntity(art))
}

func (repo *articleReaderRepository) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   uint8(art.Status),
	}
}
