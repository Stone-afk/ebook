package article

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoArticleDAO struct {
	col     *mongo.Collection
	liveCol *mongo.Collection
	node    *snowflake.Node
}

func (dao *MongoArticleDAO) GetByAuthor(ctx context.Context, authorId int64, offset, limit int) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *MongoArticleDAO) GetById(ctx context.Context, id int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *MongoArticleDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *MongoArticleDAO) ListPubByUtime(ctx context.Context, uTime time.Time, offset int, limit int) ([]PublishedArticle, error) {
	//TODO implement me
	panic("implement me")
}

func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) ArticleDAO {
	return &MongoArticleDAO{
		node:    node,
		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
	}
}

func (dao *MongoArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	// 没法子引入事务的概念
	// 首先第一步，保存制作库
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = dao.UpdateById(ctx, art)
	} else {
		id, err = dao.Insert(ctx, art)
	}
	if err != nil {
		return id, err
	}
	art.Id = id
	// 操作线上库了, upsert 语义
	now := time.Now().UnixMilli()
	//update := bson.E{"$set", art}
	//upsert := bson.E{"$setOnInsert", bson.D{bson.E{"ctime", now}}}
	art.Utime = now
	update := bson.M{
		// 更新，如果不存在，就是插入，
		"$set": PublishedArticle(art),
		// 在插入的时候，要插入 ctime
		"$setOnInsert": bson.M{"ctime": now},
	}
	filter := bson.M{"id": art.Id}
	_, err = dao.liveCol.UpdateOne(ctx, filter, update,
		options.Update().SetUpsert(true))
	return id, err
}

func (dao *MongoArticleDAO) SyncStatus(ctx context.Context, authorId, id int64, status uint8) error {
	filter := bson.D{bson.E{Key: "id", Value: id},
		bson.E{Key: "author_id", Value: authorId}}
	sets := bson.D{bson.E{Key: "$set",
		// 这里可以考虑直接使用整个 art，因为会忽略零值。
		// 参考 Sync 中的写法
		// 但是我一般都喜欢显式指定要被更新的字段，确保可读性和可维护性
		Value: bson.D{bson.E{Key: "status", Value: status}}}}
	res, err := dao.col.UpdateOne(ctx, filter, sets)
	if err != nil {
		return err
	}
	if res.ModifiedCount != 1 {
		return ErrPossibleIncorrectAuthor
	}
	return nil
}

func (dao *MongoArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	art.Id = dao.node.Generate().Int64()
	now := time.Now().UnixMilli()
	art.Utime = now
	art.Ctime = now
	_, err := dao.col.InsertOne(ctx, art)
	return art.Id, err
}

func (dao *MongoArticleDAO) UpdateById(ctx context.Context, art Article) error {
	filter := bson.D{bson.E{Key: "id", Value: art.Id},
		bson.E{Key: "author_id", Value: art.AuthorId}}
	sets := bson.D{bson.E{Key: "$set",
		// 这里可以考虑直接使用整个 art，因为会忽略零值。
		// 参考 Sync 中的写法
		// 但是我一般都喜欢显式指定要被更新的字段，确保可读性和可维护性
		Value: bson.D{bson.E{Key: "title", Value: art.Title},
			bson.E{Key: "content", Value: art.Content},
			bson.E{Key: "status", Value: art.Status},
			bson.E{Key: "utime", Value: time.Now().UnixMilli()},
		}}}
	res, err := dao.col.UpdateOne(ctx, filter, sets)
	if err != nil {
		return err
	}
	if res.MatchedCount != 1 {
		// 比较可能就是有人更新别人的文章，比如说攻击者跟你过不去
		return errors.New("更新失败")
	}
	return nil
}
