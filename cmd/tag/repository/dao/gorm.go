package dao

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
)

type GORMTagDAO struct {
	db *gorm.DB
}

func (dao *GORMTagDAO) CreateTag(ctx context.Context, tag Tag) (int64, error) {
	now := time.Now().UnixMilli()
	tag.Ctime = now
	tag.Utime = now
	err := dao.db.WithContext(ctx).Create(&tag).Error
	return tag.Id, err
}

func (dao *GORMTagDAO) CreateTagBiz(ctx context.Context, tagBiz []TagBiz) error {
	if len(tagBiz) == 0 {
		return nil
	}
	now := time.Now().UnixMilli()
	for _, t := range tagBiz {
		t.Ctime = now
		t.Utime = now
	}
	first := tagBiz[0]
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 完成了覆盖式的操作
		// 如果 tag_biz 里面没有 uid 字段。你的删除就很麻烦
		// delete from tag_biz where tid IN
		// (select distinct id from tag where uid = ?) AND biz = ? AND biz_id = ?
		err := tx.Model(&TagBiz{}).Delete(
			" uid = ? AND biz = ? AND biz_id = ?",
			first.Uid, first.BizId, first.BizId).Error
		if err != nil {
			return err
		}
		return tx.Create(&tagBiz).Error
	})
}

func (dao *GORMTagDAO) GetTagsByUid(ctx context.Context, uid int64) ([]Tag, error) {
	var res []Tag
	err := dao.db.WithContext(ctx).Where("uid= ?", uid).Find(&res).Error
	return res, err
}

func (dao *GORMTagDAO) GetTagsByBiz(ctx context.Context, uid int64, biz string, bizId int64) ([]Tag, error) {
	// 这边使用 JOIN 查询，如果你不想使用 JOIN 查询，
	// 就在 repository 里面分成两次查询
	var res []TagBiz
	err := dao.db.WithContext(ctx).Model(&TagBiz{}).
		InnerJoins("Tag", dao.db.Model(&Tag{})).
		Where("Tag.uid = ? AND biz = ? AND biz_id = ?", uid, biz, bizId).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return slice.Map(res, func(idx int, src TagBiz) Tag {
		return *src.Tag
	}), nil
	// 按照标准互联网的做法，是不用 JOIN 之类的查询，
	// 分两次
	//var tbs []TagBiz
	//err := dao.db.WithContext(ctx).Where("uid =? AND biz = ? AND biz_id = ?").Find(&tbs).Error
	//if err != nil {
	//	return nil, err
	//}
	//ids := slice.Map(tbs, func(idx int, src TagBiz) int64 {
	//	return src.Tid
	//})
	// 如果你有 id => tag 的缓存。或者 uid => tag 的缓存，你可以利用缓存
	//var res []Tag
	//err = dao.db.WithContext(ctx).Where("id IN ?", ids).Find(&res).Error

	//return res, err
}

func (dao *GORMTagDAO) GetTags(ctx context.Context, offset, limit int) ([]Tag, error) {
	var res []Tag
	err := dao.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (dao *GORMTagDAO) GetTagsById(ctx context.Context, ids []int64) ([]Tag, error) {
	var res []Tag
	err := dao.db.WithContext(ctx).Where("id IN ?", ids).Find(&res).Error
	return res, err
}

func NewGORMTagDAO(db *gorm.DB) TagDAO {
	return &GORMTagDAO{
		db: db,
	}
}
