package repository

import (
	"context"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/repository/dao"
)

type tagRepository struct {
	dao dao.TagDAO
}

func (repo *tagRepository) InputTag(ctx context.Context, msg domain.Tag) error {
	//TODO implement me
	panic("implement me")
}

func (repo *tagRepository) SearchTag(ctx context.Context, uid int64, biz string, keywords []string) ([]domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func NewTagRepository(d dao.TagDAO) TagRepository {
	return &tagRepository{
		dao: d,
	}
}
