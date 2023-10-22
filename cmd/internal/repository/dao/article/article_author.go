package article

import (
	"context"
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

// UpdateById 只更新标题和
func (dao *GORMArticleAuthorDAO) UpdateById(ctx context.Context,
	art Article) error {
	panic("")
}
