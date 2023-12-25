package dao

import (
	"context"
	"ebook/cmd/internal/repository/dao/article"
	"ebook/cmd/internal/repository/dao/async_sms"
	"ebook/cmd/internal/repository/dao/job"
	"ebook/cmd/internal/repository/dao/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"time"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&user.User{},
		&article.Article{},
		&article.PublishedArticle{},
		&article.PublishedArticleV1{},
		&async_sms.AsyncSms{},
		&job.Job{},
	)
}

func InitCollections(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	index := []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "author_id", Value: 1},
				bson.E{Key: "ctime", Value: 1},
			},
			Options: options.Index(),
		},
	}
	_, err := db.Collection("articles").
		Indexes().CreateMany(ctx, index)
	if err != nil {
		return err
	}
	_, err = db.Collection("published_articles").
		Indexes().CreateMany(ctx, index)
	return err
}
