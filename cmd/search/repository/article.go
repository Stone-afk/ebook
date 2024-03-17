package repository

import (
	"context"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/repository/dao"
	"github.com/ecodeclub/ekit/slice"
)

type articleRepository struct {
	dao    dao.ArticleDAO
	tagDao dao.TagDAO
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
	tags, err := repo.tagDao.Search(ctx, uid, "article", keywords)
	if err != nil {
		return nil, err
	}
	ids := slice.Map(tags, func(idx int, src dao.BizTags) int64 {
		return src.BizId
	})
	arts, err := repo.dao.Search(ctx, ids, keywords)
	if err != nil {
		return nil, err
	}
	return slice.Map(arts, func(idx int, src dao.Article) domain.Article {
		return domain.Article{
			Id:      src.Id,
			Title:   src.Title,
			Status:  src.Status,
			Content: src.Content,
			Tags:    src.Tags,
		}
	}), nil
}

func NewArticleRepository(dao dao.ArticleDAO, tagDao dao.TagDAO) ArticleRepository {
	return &articleRepository{
		dao:    dao,
		tagDao: tagDao,
	}
}
