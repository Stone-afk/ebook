package repository

import (
	"context"
	"ebook/cmd/followrelation/domain"
	"ebook/cmd/followrelation/repository/cache"
	"ebook/cmd/followrelation/repository/dao"
	"ebook/cmd/pkg/logger"
)

type CachedRelationRepository struct {
	dao   dao.FollowRelationDao
	cache cache.FollowCache
	l     logger.Logger
}

func (repo *CachedRelationRepository) GetFollower(ctx context.Context, followee, offset, limit int64) ([]domain.FollowRelation, error) {
	followerList, err := repo.dao.FindFollowerList(ctx, followee, offset, limit)
	if err != nil {
		return nil, err
	}
	return repo.genFollowRelationList(followerList), nil
}

func (repo *CachedRelationRepository) GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error) {
	followeeList, err := repo.dao.FindFolloweeList(ctx, follower, offset, limit)
	if err != nil {
		return nil, err
	}
	return repo.genFollowRelationList(followeeList), nil
}

func (repo *CachedRelationRepository) genFollowRelationList(followerList []dao.FollowRelation) []domain.FollowRelation {
	res := make([]domain.FollowRelation, 0, len(followerList))
	for _, c := range followerList {
		res = append(res, repo.toDomain(c))
	}
	return res
}

func (repo *CachedRelationRepository) FollowInfo(ctx context.Context, follower int64, followee int64) (domain.FollowRelation, error) {
	c, err := repo.dao.FollowRelationDetail(ctx, follower, followee)
	if err != nil {
		return domain.FollowRelation{}, err
	}
	return repo.toDomain(c), nil
}

func (repo *CachedRelationRepository) AddFollowRelation(ctx context.Context, f domain.FollowRelation) error {
	err := repo.dao.CreateFollowRelation(ctx, repo.toEntity(f))
	if err != nil {
		return err
	}
	return repo.cache.Follow(ctx, f.Follower, f.Followee)
}

func (repo *CachedRelationRepository) InactiveFollowRelation(ctx context.Context, follower int64, followee int64) error {
	err := repo.dao.UpdateStatus(ctx, followee, follower, dao.FollowRelationStatusInactive)
	if err != nil {
		return err
	}
	return repo.cache.CancelFollow(ctx, follower, followee)
}

func (repo *CachedRelationRepository) GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	// 快路径
	res, err := repo.cache.StaticsInfo(ctx, uid)
	if err == nil {
		return res, err
	}
	// 慢路径
	res.Followers, err = repo.dao.CntFollower(ctx, uid)
	if err != nil {
		return res, err
	}
	res.Followees, err = repo.dao.CntFollowee(ctx, uid)
	if err != nil {
		return res, err
	}
	err = repo.cache.SetStaticsInfo(ctx, uid, res)
	if err != nil {
		// 这里记录日志
		repo.l.Error("缓存关注统计信息失败",
			logger.Error(err),
			logger.Int64("uid", uid))
	}
	return res, nil
}

func (repo *CachedRelationRepository) toDomain(fr dao.FollowRelation) domain.FollowRelation {
	return domain.FollowRelation{
		Id:       fr.ID,
		Followee: fr.Followee,
		Follower: fr.Follower,
	}
}

func (repo *CachedRelationRepository) toEntity(c domain.FollowRelation) dao.FollowRelation {
	return dao.FollowRelation{
		Followee: c.Followee,
		Follower: c.Follower,
	}
}

func NewFollowRelationRepository(dao dao.FollowRelationDao,
	cache cache.FollowCache, l logger.Logger) FollowRepository {
	return &CachedRelationRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}
