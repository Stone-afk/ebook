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

// PreloadUserTags 在 toB 的场景下，你可以提前预加载缓存
func (repo *CachedTagRepository) PreloadUserTags(ctx context.Context) error {
	// 怎么预加载？
	// 缓存里面，究竟怎么存？
	// 1. 放 json，json 里面是一个用户的所有的标签 uid => [{}, {}]
	// 按照用户 ID 来查找
	//var uid int64= 1
	//for {
	//	repo.dao.GetTagsByUid(ctx, uid)
	//	uid ++
	//}
	// select * from tags group by uid
	// 使用 redis 的数据结构
	// 1. list  使用 List 在多个实例启动时预加载数据可能导致冲突，所以改成 hash
	// 2. hash 用 hash 结构
	// 3. set, sorted set 都可以
	//TODO implement me
	panic("implement me")
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
