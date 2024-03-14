package repository

import (
	"context"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/repository/dao"
)

type articleRepository struct {
	dao dao.ArticleDAO
}

func (repo *articleRepository) InputArticle(ctx context.Context, msg domain.Article) error {
	return repo.dao.InputArticle(ctx, dao.Article{
		Id:      msg.Id,
		Title:   msg.Title,
		Status:  msg.Status,
		Content: msg.Content,
	})
}

func (repo *articleRepository) SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func NewArticleRepository(d dao.ArticleDAO) ArticleRepository {
	return &articleRepository{
		dao: d,
	}
}
