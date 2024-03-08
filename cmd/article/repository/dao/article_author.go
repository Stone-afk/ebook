package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// GORMArticleAuthorDAO 复制了 GORMArticleDAO 的代码
type GORMArticleAuthorDAO struct {
	db *gorm.DB
}

func (dao *GORMArticleAuthorDAO) Create(ctx context.Context,
	art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

// UpdateById 只更新标题和内容
func (dao *GORMArticleAuthorDAO) UpdateById(ctx context.Context,
	art Article) error {
	now := time.Now().UnixMilli()
	// 依赖 gorm 忽略零值的特性，会用主键进行更新
	// 可读性很差
	res := dao.db.WithContext(ctx).Model(&art).
		Where("id=? AND author_id = ?", art.Id, art.AuthorId).
		// 当你用这种每次都指定被更新列的写法
		// 可读性强，但是每一次更新更多的列的时候，你都要修改
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   now,
		})
	// 你要不要检查真的更新了没？
	// res.RowsAffected // 更新行数
	err := res.Error
	if err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		//dangerousDBOp.Count(1)
		// 补充一点日志
		return fmt.Errorf("更新失败，可能是创作者非法 id %d, author_id %d",
			art.Id, art.AuthorId)
	}
	return nil
}
