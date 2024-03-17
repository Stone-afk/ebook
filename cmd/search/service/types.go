package service

import (
	"context"
	"ebook/cmd/search/domain"
)

type SearchService interface {
	Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error)
}

type TagService interface {
	Search(ctx context.Context, uid int64, biz string, expression string) (domain.SearchResult, error)
}

type ArticleSearchService interface {
	SearchArticle(ctx context.Context, uid int64, expression string) (domain.SearchResult, error)
}

type UserSearchService interface {
	SearchUser(ctx context.Context, expression string) (domain.SearchResult, error)
}

type SyncService interface {
	InputTag(ctx context.Context, tag domain.BizTags) error
	InputArticle(ctx context.Context, article domain.Article) error
	InputUser(ctx context.Context, user domain.User) error
	InputAny(ctx context.Context, index, docID, data string) error
}
