package repository

import (
	"context"
	"ebook/cmd/internal/repository/cache"
	"ebook/cmd/internal/repository/dao/interactive"
	"ebook/cmd/pkg/logger"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/interactive.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/mocks/interactive.mock.go
type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
}

type interactiveRepository struct {
	cache cache.InteractiveCache
	dao   interactive.InteractiveDAO
	l     logger.Logger
}

func NewCachedInteractiveRepository(dao interactive.InteractiveDAO,
	cache cache.InteractiveCache, l logger.Logger) InteractiveRepository {
	return &interactiveRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (repo *interactiveRepository) IncrLike(ctx context.Context,
	biz string, bizId int64, uid int64) error {
	// 先插入点赞，然后更新点赞计数，更新缓存
	err := repo.dao.InsertLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return repo.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}

func (repo *interactiveRepository) DecrLike(ctx context.Context,
	biz string, bizId, uid int64) error {
	err := repo.dao.DeleteLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return repo.cache.DecrLikeCntIfPresent(ctx, biz, bizId)
}

func (repo *interactiveRepository) IncrReadCnt(ctx context.Context,
	biz string, bizId int64) error {
	// 要考虑缓存方案了
	// 这两个操作能不能换顺序？ —— 不能
	err := repo.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	//go func() {
	//	c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
	//}()
	//return err
	return repo.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}