package dao

import "context"

type TagDAO interface {
	CreateTag(ctx context.Context, tag Tag) (int64, error)
	CreateTagBiz(ctx context.Context, tagBiz []TagBiz) error
	GetTagsByUid(ctx context.Context, uid int64) ([]Tag, error)
	GetTagsByBiz(ctx context.Context, uid int64, biz string, bizId int64) ([]Tag, error)
	GetTags(ctx context.Context, offset, limit int) ([]Tag, error)
	GetTagsById(ctx context.Context, ids []int64) ([]Tag, error)
}
