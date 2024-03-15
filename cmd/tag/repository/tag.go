package repository

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/tag/domain"
	"ebook/cmd/tag/repository/cache"
	"ebook/cmd/tag/repository/dao"
	"github.com/ecodeclub/ekit/slice"
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
	id, err := repo.dao.CreateTag(ctx, repo.toEntity(tag))
	if err != nil {
		return 0, err
	}
	// 也可以考虑用 DelTags
	// 记得更新你的缓存
	err = repo.cache.Append(ctx, tag.Uid, tag)
	if err != nil {
		// 记录日志
		repo.l.Error("缓存自定义标签失败",
			logger.Error(err),
			logger.Int64("uid", tag.Uid))
	}
	return id, nil
}

func (repo *CachedTagRepository) BindTagToBiz(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error {
	// 按照需求，是要覆盖式地执行打标签
	// 新的标签完全覆盖老的标签
	// 按道理应该是 repository 去控制的
	return repo.dao.CreateTagBiz(ctx, slice.Map(tags, func(idx int, src int64) dao.TagBiz {
		return dao.TagBiz{
			Tid:   src,
			BizId: bizId,
			Biz:   biz,
			Uid:   uid,
		}
	}))
}

func (repo *CachedTagRepository) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	res, err := repo.cache.GetTags(ctx, uid)
	if err == nil {
		return res, nil
	}
	// 下面也是慢路径，你同样可以说降级的时候不执行

	// 如果我要缓存
	// 我这里应该是 uid => tags 的映射
	// toB 的时候，我直接全量缓存
	// 我要在应用启动的时候，把缓存加载好
	// 如果你认为你的 tags 是没有过期时间的，你这里就不用回查数据库了
	tags, err := repo.dao.GetTagsByUid(ctx, uid)
	if err != nil {
		return nil, err
	}
	res = slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return repo.toDomain(src)
	})
	err = repo.cache.Append(ctx, uid, res...)
	if err != nil {
		// 记录日志
		// 缓存回写失败，不认为是一个问题
		repo.l.Error("缓存自定义标签列表失败",
			logger.Error(err),
			logger.Int64("uid", uid))
	}
	return res, nil
}

func (repo *CachedTagRepository) GetTagsById(ctx context.Context, ids []int64) ([]domain.Tag, error) {
	tags, err := repo.dao.GetTagsById(ctx, ids)
	if err != nil {
		return nil, err
	}
	return slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return repo.toDomain(src)
	}), nil
}

func (repo *CachedTagRepository) GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error) {
	// 要缓存的话，就是 uid + biz + biz_id 构成一个 key
	tags, err := repo.dao.GetTagsByBiz(ctx, uid, biz, bizId)
	if err != nil {
		return nil, err
	}
	return slice.Map(tags, func(idx int, src dao.Tag) domain.Tag {
		return repo.toDomain(src)
	}), nil
}

func (repo *CachedTagRepository) toDomain(tag dao.Tag) domain.Tag {
	return domain.Tag{
		Id:   tag.Id,
		Name: tag.Name,
		Uid:  tag.Uid,
	}
}

func (repo *CachedTagRepository) toEntity(tag domain.Tag) dao.Tag {
	return dao.Tag{
		Id:   tag.Id,
		Name: tag.Name,
		Uid:  tag.Uid,
	}
}

func NewTagRepository(tagDAO dao.TagDAO, c cache.TagCache, l logger.Logger) TagRepository {
	return &CachedTagRepository{
		dao:   tagDAO,
		l:     l,
		cache: c,
	}
}
