package repository

import (
	"context"
	"ebook/cmd/article/domain"
	"ebook/cmd/article/repository/dao"
)

type articleReaderRepository struct {
	dao dao.ArticleReaderDAO
}

func NewArticleReaderRepository(dao dao.ArticleReaderDAO) ArticleReaderRepository {
	return &articleReaderRepository{
		dao: dao,
	}
}

func (repo *articleReaderRepository) Save(ctx context.Context, art domain.Article) error {
	return repo.dao.Upsert(ctx, repo.toEntity(art))
}

func (repo *articleReaderRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   uint8(art.Status),
	}
}
