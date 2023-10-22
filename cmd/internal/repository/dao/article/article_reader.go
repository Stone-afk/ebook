package article

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GORMArticleReaderDAO struct {
	db *gorm.DB
}

func NewGORMArticleReaderDAO(db *gorm.DB) ArticleReaderDAO {
	return &GORMArticleReaderDAO{
		db: db,
	}
}

func (dao *GORMArticleReaderDAO) Upsert(ctx context.Context, art Article) error {
	return dao.db.Clauses(clause.OnConflict{
		// ID 冲突的时候。实际上，在 MYSQL 里面你写不写都可以
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
		}),
	}).Create(&art).Error
}
