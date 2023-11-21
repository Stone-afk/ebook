package repository

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository/cache"
	"ebook/cmd/internal/repository/dao/interactive"
	"ebook/cmd/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/interactive.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/mocks/interactive.mock.go
type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	BatchIncrReadCnt(ctx context.Context, biz []string, bizIds []int64) error
	IncrLike(ctx context.Context, biz string, bizId, userId int64) error
	DecrLike(ctx context.Context, biz string, bizId, userId int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, userId int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, userId int64) (bool, error)
	AddRecord(ctx context.Context, aid int64, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error
	GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error)
}

type interactiveRepository struct {
	cache cache.InteractiveCache
	dao   interactive.InteractiveDAO
	l     logger.Logger
}

func NewInteractiveRepository(dao interactive.InteractiveDAO,
	cache cache.InteractiveCache, l logger.Logger) InteractiveRepository {
	return &interactiveRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}

func (repo *interactiveRepository) GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error) {
	vals, err := repo.dao.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}
	return slice.Map[interactive.Interactive, domain.Interactive](vals,
		func(idx int, src interactive.Interactive) domain.Interactive {
			return repo.toDomain(src)
		}), nil
}

func (repo *interactiveRepository) AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error {
	err := repo.dao.InsertCollectionBiz(ctx, interactive.UserCollectionBiz{
		Biz:   biz,
		BizId: bizId,
		Cid:   cid,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	return repo.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
}

func (repo *interactiveRepository) AddRecord(ctx context.Context, aid int64, uid int64) error {
	//TODO implement me
	panic("implement me")
}

// BatchIncrReadCnt bizs 和 ids 的长度必须相等
func (repo *interactiveRepository) BatchIncrReadCnt(ctx context.Context,
	bizs []string, bizIds []int64) error {
	// 在这里要不要检测 bizs 和 ids 的长度是否相等？
	err := repo.dao.BatchIncrReadCnt(ctx, bizs, bizIds)
	if err != nil {
		return err
	}
	// 你也要批量的去修改 redis，所以就要去改 lua 脚本
	// c.cache.IncrReadCntIfPresent()
	// TODO, 等写新的 lua 脚本/或者用 pipeline
	return nil
}

func (repo *interactiveRepository) Get(ctx context.Context,
	biz string, bizId int64) (domain.Interactive, error) {
	// 要从缓存拿出来阅读数，点赞数和收藏数
	intr, err := repo.cache.Get(ctx, biz, bizId)
	if err == nil {
		return intr, nil
	}
	// 但不是所有的结构体都是可比较的
	//if intr == (domain.Interactive{}) {
	//
	//}
	// 在这里查询数据库
	daoIntr, err := repo.dao.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	intr = repo.toDomain(daoIntr)
	go func() {
		er := repo.cache.Set(ctx, biz, bizId, intr)
		// 记录日志
		if er != nil {
			repo.l.Error("回写缓存失败",
				logger.String("biz", biz),
				logger.Int64("bizId", bizId),
			)
		}
	}()
	return intr, nil
}

func (repo *interactiveRepository) Liked(ctx context.Context,
	biz string, id int64, userId int64) (bool, error) {
	_, err := repo.dao.GetLikeInfo(ctx, biz, id, userId)
	switch err {
	case nil:
		return true, nil
	case interactive.ErrRecordNotFound:
		// 你要吞掉
		return false, nil
	default:
		return false, err
	}
}

func (repo *interactiveRepository) Collected(ctx context.Context,
	biz string, id int64, userId int64) (bool, error) {
	_, err := repo.dao.GetCollectionInfo(ctx, biz, id, userId)
	switch err {
	case nil:
		return true, nil
	case interactive.ErrRecordNotFound:
		// 你要吞掉
		return false, nil
	default:
		return false, err
	}
}

func (repo *interactiveRepository) IncrLike(ctx context.Context,
	biz string, bizId int64, userId int64) error {
	// 先插入点赞，然后更新点赞计数，更新缓存
	err := repo.dao.InsertLikeInfo(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}
	return repo.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}

func (repo *interactiveRepository) DecrLike(ctx context.Context,
	biz string, bizId, userId int64) error {
	err := repo.dao.DeleteLikeInfo(ctx, biz, bizId, userId)
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

func (repo *interactiveRepository) toDomain(intr interactive.Interactive) domain.Interactive {
	return domain.Interactive{
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		ReadCnt:    intr.ReadCnt,
	}
}
