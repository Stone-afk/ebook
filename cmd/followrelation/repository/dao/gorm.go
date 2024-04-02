package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GORMFollowRelationDAO struct {
	db *gorm.DB
}

func (dao *GORMFollowRelationDAO) FindFollowerList(ctx context.Context, followee, offset, limit int64) ([]FollowRelation, error) {
	var res []FollowRelation
	err := dao.db.WithContext(ctx).
		Where("followee = ? AND status = ?", followee, FollowRelationStatusActive).
		Offset(int(offset)).Limit(int(limit)).
		Find(&res).Error
	return res, err
}

func (dao *GORMFollowRelationDAO) FindFolloweeList(ctx context.Context, follower, offset, limit int64) ([]FollowRelation, error) {
	var res []FollowRelation
	err := dao.db.WithContext(ctx).
		Where("follower = ? AND status = ?", follower, FollowRelationStatusActive).
		Offset(int(offset)).Limit(int(limit)).
		Find(&res).Error
	return res, err
}

func (dao *GORMFollowRelationDAO) FollowRelationDetail(ctx context.Context, follower int64, followee int64) (FollowRelation, error) {
	var res FollowRelation
	err := dao.db.WithContext(ctx).Where("follower = ? AND followee = ? AND status = ?",
		follower, followee, FollowRelationStatusActive).First(&res).Error
	return res, err
}

func (dao *GORMFollowRelationDAO) CreateFollowRelation(ctx context.Context, f FollowRelation) error {
	now := time.Now().UnixMilli()
	f.Utime = now
	f.Ctime = now
	f.Status = FollowRelationStatusActive
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"utime":  now,
			"status": FollowRelationStatusActive,
		}),
	}).Create(&f).Error
}

func (dao *GORMFollowRelationDAO) UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).
		Where("follower = ? AND followee = ?", follower, followee).
		Updates(map[string]any{
			"status": status,
			"utime":  now,
		}).Error
}

// CntFollower 统计粉丝数量
func (dao *GORMFollowRelationDAO) CntFollower(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := dao.db.WithContext(ctx).
		Select("count(follower)").
		Where("followee = ? AND status = ?",
			uid, FollowRelationStatusActive).Count(&res).Error
	return res, err
}

// CntFollowee 统计关注数量
func (dao *GORMFollowRelationDAO) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	var res int64
	err := dao.db.WithContext(ctx).
		Select("count(followee)").
		Where("follower = ? AND status = ?",
			uid, FollowRelationStatusActive).Count(&res).Error
	return res, err
}

func NewGORMFollowRelationDAO(db *gorm.DB) FollowRelationDao {
	return &GORMFollowRelationDAO{
		db: db,
	}
}
