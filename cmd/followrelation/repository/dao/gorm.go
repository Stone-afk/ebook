package dao

import (
	"context"
	"gorm.io/gorm"
)

type GORMFollowRelationDAO struct {
	db *gorm.DB
}

func (dao *GORMFollowRelationDAO) FollowRelationList(ctx context.Context, follower, offset, limit int64) ([]FollowRelation, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMFollowRelationDAO) FollowRelationDetail(ctx context.Context, follower int64, followee int64) (FollowRelation, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMFollowRelationDAO) CreateFollowRelation(ctx context.Context, c FollowRelation) error {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMFollowRelationDAO) UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMFollowRelationDAO) CntFollower(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := dao.db.WithContext(ctx).
		Select("count(follower)").
		Where("followee = ? AND status = ?",
			uid, FollowRelationStatusActive).Count(&res).Error
	return res, err
}

func (dao *GORMFollowRelationDAO) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewGORMFollowRelationDAO(db *gorm.DB) FollowRelationDao {
	return &GORMFollowRelationDAO{
		db: db,
	}
}
