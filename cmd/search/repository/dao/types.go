package dao

import "context"

type TagDAO interface {
	InputTag(ctx context.Context, tag Tag) error
	Search(ctx context.Context, uid int64, biz string, keywords []string) ([]Tag, error)
}

type UserDAO interface {
	InputUser(ctx context.Context, user User) error
	Search(ctx context.Context, keywords []string) ([]User, error)
}

type ArticleDAO interface {
	InputArticle(ctx context.Context, article Article) error
	Search(ctx context.Context, tagArtIds []int64, keywords []string) ([]Article, error)
}

type AnyDAO interface {
	Input(ctx context.Context, index, docID, data string) error
}
