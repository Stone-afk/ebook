package repository

import (
	"context"
	"ebook/cmd/search/domain"
)

type AnyRepository interface {
	Input(ctx context.Context, index string, docID string, data string) error
}

type UserRepository interface {
	InputUser(ctx context.Context, msg domain.User) error
	SearchUser(ctx context.Context, keywords []string) ([]domain.User, error)
}

type ArticleRepository interface {
	InputArticle(ctx context.Context, msg domain.Article) error
	SearchArticle(ctx context.Context, uid int64, keywords []string) ([]domain.Article, error)
}

type TagRepository interface {
	InputTag(ctx context.Context, msg domain.BizTags) error
	SearchTag(ctx context.Context, uid int64, biz string, keywords []string) ([]domain.BizTags, error)
}
