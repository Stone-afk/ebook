package article

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoArticleDAO struct {
	col     *mongo.Collection
	liveCol *mongo.Collection
	node    *snowflake.Node
}

func (dao *MongoArticleDAO) Sync(ctx context.Context, art Article) (int64, error) {
	panic("")
}

func (dao *MongoArticleDAO) SyncStatus(ctx context.Context, authorId, id int64, status uint8) error {
	panic("")
}

func (dao *MongoArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	panic("")
}

func (dao *MongoArticleDAO) UpdateById(ctx context.Context, art Article) error {
	panic("")
}
