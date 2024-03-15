package dao

import (
	"context"
	"gorm.io/gorm"
)

type GORMTagDAO struct {
	db *gorm.DB
}

func (dao *GORMTagDAO) CreateTag(ctx context.Context, tag Tag) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMTagDAO) CreateTagBiz(ctx context.Context, tagBiz []TagBiz) error {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMTagDAO) GetTagsByUid(ctx context.Context, uid int64) ([]Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMTagDAO) GetTagsByBiz(ctx context.Context, uid int64, biz string, bizId int64) ([]Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMTagDAO) GetTags(ctx context.Context, offset, limit int) ([]Tag, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *GORMTagDAO) GetTagsById(ctx context.Context, ids []int64) ([]Tag, error) {
	//TODO implement me
	panic("implement me")
}

func NewGORMTagDAO(db *gorm.DB) TagDAO {
	return &GORMTagDAO{
		db: db,
	}
}
