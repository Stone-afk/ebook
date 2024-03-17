package service

import (
	"context"
	"ebook/cmd/search/domain"
)

type SearchService interface {
	Search(ctx context.Context, uid int64, expression string) (domain.SearchResult, error)
}

type TagService interface {
	SearchBizTags(ctx context.Context, uid int64, biz string, expression string) (domain.SearchBizTagsResult, error)
}

type ArticleSearchService interface {
	SearchArticle(ctx context.Context, uid int64, expression string) (domain.SearchArticleResult, error)
}

type UserSearchService interface {
	SearchUser(ctx context.Context, expression string) (domain.SearchUserResult, error)
}

type SyncService interface {
	InputBizTags(ctx context.Context, tag domain.BizTags) error
	InputArticle(ctx context.Context, article domain.Article) error
	InputUser(ctx context.Context, user domain.User) error
	InputAny(ctx context.Context, index, docID, data string) error
}
