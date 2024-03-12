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

func (repo *CachedRelationRepository) GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *CachedRelationRepository) FollowInfo(ctx context.Context, follower int64, followee int64) (domain.FollowRelation, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *CachedRelationRepository) AddFollowRelation(ctx context.Context, f domain.FollowRelation) error {
	//TODO implement me
	panic("implement me")
}

func (repo *CachedRelationRepository) InactiveFollowRelation(ctx context.Context, follower int64, followee int64) error {
	//TODO implement me
	panic("implement me")
}

func (repo *CachedRelationRepository) GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *CachedRelationRepository) toDomain(fr dao.FollowRelation) domain.FollowRelation {
	return domain.FollowRelation{
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
