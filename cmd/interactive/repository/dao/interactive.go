package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GORMInteractiveDAO struct {
	db *gorm.DB
}

func NewGORMInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{
		db: db,
	}
}

func (dao *GORMInteractiveDAO) GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error) {
	var res []Interactive
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND id IN ?", biz, ids).Find(&res).Error
	return res, err
}

// BatchIncrReadCnt
// biz = a, bizid = 1
// biz = a, bizid = 1
// biz = a, bizid = 1
// biz = a bizId = 2
// biz = a bizId = 2
// biz = a bizId = 2
// biz = a bizId = 2
func (dao *GORMInteractiveDAO) BatchIncrReadCnt(ctx context.Context,
	bizs []string, ids []int64) error {
	// 可以用 map 合并吗？
	// 看情况。如果一批次里面，biz 和 bizid 都相等的占很多，那么就map 合并，性能会更好
	// 不然你合并了没有效果

	// 为什么快？
	// A：十条消息调用十次 IncrReadCnt，
	// B 就是批量
	// 事务本身的开销，A 是 B 的十倍
	// 刷新 redolog, undolog, binlog 到磁盘，A 是十次，B 是一次
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMInteractiveDAO(tx)
		for i := range bizs {
			err := txDAO.IncrReadCnt(ctx, bizs[i], ids[i])
			if err != nil {
				// 记个日志就拉到
				// 也可以 return err
				return err
			}
		}
		return nil
	})
}

func (dao *GORMInteractiveDAO) Get(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	var res Interactive
	err := dao.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ?", biz, bizId).
		First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) GetLikeInfo(ctx context.Context,
	biz string, bizId, userId int64) (UserLikeBiz, error) {
	var res UserLikeBiz
	err := dao.db.WithContext(ctx).
		Where("biz=? AND biz_id = ? AND uid = ? AND status = ?",
			biz, bizId, userId, 1).First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) GetCollectionInfo(ctx context.Context,
	biz string, bizId, userId int64) (UserCollectionBiz, error) {
	var res UserCollectionBiz
	err := dao.db.WithContext(ctx).
		Where("biz=? AND biz_id = ? AND uid = ?", biz, bizId, userId).
		First(&res).Error
	return res, err
}

func (dao *GORMInteractiveDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	// 一把梭
	// 同时记录点赞，以及更新点赞计数
	// 首先你需要一张表来记录，谁点给什么资源点了赞
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.WithContext(ctx).Create(&cb).Error
		if err != nil {
			return err
		}
		// 这边就是更新数量
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"collect_cnt": gorm.Expr("`collect_cnt`+1"),
				"utime":       now,
			}),
		}).Create(&Interactive{
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
			Biz:        cb.Biz,
			BizId:      cb.BizId,
		}).Error
	})
}

func (dao *GORMInteractiveDAO) InsertLikeInfo(ctx context.Context, biz string, bizId, userId int64) error {
	// 一把梭
	// 同时记录点赞，以及更新点赞计数
	// 首先你需要一张表来记录，谁点给什么资源点了赞
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先准备插入点赞记录
		// 有没有可能已经点赞过了？
		// 我要不要校验一下，这里必须是没有点赞过
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"utime":  now,
				"status": 1,
			}),
		}).Create(&UserLikeBiz{
			Biz:    biz,
			BizId:  bizId,
			Uid:    userId,
			Status: 1,
			Ctime:  now,
			Utime:  now,
		}).Error
		if err != nil {
			return err
		}
		return tx.Clauses(clause.OnConflict{
			// MySQL 不写
			//Columns:
			DoUpdates: clause.Assignments(map[string]any{
				"like_cnt": gorm.Expr("like_cnt + 1"),
				"utime":    time.Now().UnixMilli(),
			}),
		}).Create(&Interactive{
			Biz:     biz,
			BizId:   bizId,
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
		}).Error
	})
}

func (dao *GORMInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId, userId int64) error {
	now := time.Now().UnixMilli()
	// 控制事务超时
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 两个操作
		// 一个是软删除点赞记录
		// 一个是减点赞数量
		err := tx.Model(&UserLikeBiz{}).
			Where("biz=? AND biz_id = ? AND userId = ?", biz, bizId, userId).
			Updates(map[string]any{
				"utime":  now,
				"status": 0,
			}).Error
		if err != nil {
			return err
		}
		return tx.Model(&Interactive{}).
			Where("biz=? AND biz_id = ?", biz, bizId).
			Updates(map[string]any{
				"utime":    now,
				"like_cnt": gorm.Expr("like_cnt-1"),
			}).Error
	})
}

// IncrReadCnt 是一个插入或者更新语义
func (dao *GORMInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	// DAO 要怎么实现？表结构该怎么设计？
	//var intr Interactive
	//err := dao.db.
	//	Where("biz_id =? AND biz = ?", bizId, biz).
	//	First(&intr).Error
	// 两个 goroutine 过来，你查询到 read_cnt 都是 10
	//if err != nil {
	//	return err
	//}
	// 都变成了 11
	//cnt := intr.ReadCnt + 1
	//// 最终变成 11
	//dao.db.Where("biz_id =? AND biz = ?", bizId, biz).Updates(map[string]any{
	//	"read_cnt": cnt,
	//})

	// update a = a + 1
	// 数据库帮你解决并发问题
	// 有一个没考虑到，就是，我可能根本没这一行
	// 事实上这里是一个 upsert 的语义
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Clauses(clause.OnConflict{
		// MySQL 不写
		//Columns:
		DoUpdates: clause.Assignments(map[string]any{
			"read_cnt": gorm.Expr("read_cnt + 1"),
			"utime":    now,
		}),
	}).Create(&Interactive{
		Biz:     biz,
		BizId:   bizId,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}
