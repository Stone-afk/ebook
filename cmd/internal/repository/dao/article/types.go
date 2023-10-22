package article

import (
	"errors"
	"golang.org/x/net/context"
)

var ErrPossibleIncorrectAuthor = errors.New("用户在尝试操作非本人数据")

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	SyncStatus(ctx context.Context, authorId, id int64, status uint8) error
}

type ArticleReaderDAO interface {
	// Upsert INSERT OR UPDATE 语义，一般简写为 Upsert
	// 将会更新标题和内容，但是不会更新别的内容
	// 这个要求 Reader 和 Author 是不同库
	Upsert(ctx context.Context, art Article) error
	// UpsertV2 版本用于同库不同表
	UpsertV2(ctx context.Context, art PublishedArticle) error
}

// ArticleAuthorDAO 用于演示代表
type ArticleAuthorDAO interface {
	Create(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
}
