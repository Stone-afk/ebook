package repository

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/tag/domain"
	"ebook/cmd/tag/repository/cache"
	"ebook/cmd/tag/repository/dao"
)

type CachedTagRepository struct {
	dao   dao.TagDAO
	cache cache.TagCache
	l     logger.Logger
}

func (repo *CachedTagRepository) CreateTag(ctx context.Context, tag domain.Tag) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *CachedTagRepository) BindTagToBiz(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error {
	//TODO implement me
	panic("implement me")
}

func (repo *CachedTagRepository) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *CachedTagRepository) GetTagsById(ctx context.Context, ids []int64) ([]domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *CachedTagRepository) GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error) {
	//TODO implement me
	panic("implement me")
}

func NewTagRepository(tagDAO dao.TagDAO, c cache.TagCache, l logger.Logger) TagRepository {
	return &CachedTagRepository{
		dao:   tagDAO,
		l:     l,
		cache: c,
	}
}
