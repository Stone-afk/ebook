package repository

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository/dao/article"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/article_author.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/mocks/article_author.mock.go

// ArticleAuthorRepository 演示在 service 层面上分流
type ArticleAuthorRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
}

// articleAuthorRepository 按照道理，这里也是可以搞缓存的
type articleAuthorRepository struct {
	dao article.ArticleDAO
}

func NewArticleAuthorRepository(dao article.ArticleDAO) ArticleAuthorRepository {
	return &articleAuthorRepository{
		dao: dao,
	}
}

func (repo *articleAuthorRepository) Update(ctx context.Context, art domain.Article) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(art))
}

func (repo *articleAuthorRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(art))
}

func (repo *articleAuthorRepository) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   uint8(art.Status),
	}
}
