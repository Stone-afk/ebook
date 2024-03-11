package dao

import (
	"context"
	"gorm.io/gorm"
)

type GORMCommentDAO struct {
	db *gorm.DB
}

func (dao *GORMCommentDAO) Insert(ctx context.Context, u Comment) error {
	return dao.db.WithContext(ctx).Create(u).Error
}

func (dao *GORMCommentDAO) FindByBiz(ctx context.Context, biz string, bizId, minID, limit int64) ([]Comment, error) {
	var res []Comment
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND id < ? AND pid IS NULL", biz, bizId, minID).
		Limit(int(limit)).
		Find(&res).Error
	return res, err
}

func (dao *GORMCommentDAO) FindCommentList(ctx context.Context, u Comment) ([]Comment, error) {
	var res []Comment
	builder := dao.db.WithContext(ctx)
	if u.Id == 0 {
		builder = builder.
			Where("biz=?", u.Biz).
			Where("biz_id=?", u.BizID).
			Where("root_id is null")
	} else {
		builder = builder.Where("root_id=? or id =?", u.Id, u.Id)
	}
	err := builder.Find(&res).Error
	return res, err
}

func (dao *GORMCommentDAO) FindRepliesByPid(ctx context.Context, pid int64, offset, limit int) ([]Comment, error) {
	var res []Comment
	err := dao.db.WithContext(ctx).Where("pid = ?", pid).
		Order("id DESC").Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (dao *GORMCommentDAO) Delete(ctx context.Context, u Comment) error {
	return dao.db.WithContext(ctx).Delete(&Comment{Id: u.Id}).Error
}

func (dao *GORMCommentDAO) FindOneByIDs(ctx context.Context, ids []int64) ([]Comment, error) {
	var res []Comment
	err := dao.db.WithContext(ctx).
		Where("id in ?", ids).
		First(&res).
		Error
	return res, err
}

func (dao *GORMCommentDAO) FindRepliesByRid(ctx context.Context, rid int64, maxID int64, limit int64) ([]Comment, error) {
	var res []Comment
	err := dao.db.WithContext(ctx).
		Where("root_id = ? AND id > ?", rid, maxID).
		Order("id ASC").Limit(int(limit)).Find(&res).Error
	return res, err
}

func NewCommentDAO(db *gorm.DB) CommentDAO {
	return &GORMCommentDAO{
		db: db,
	}
}
