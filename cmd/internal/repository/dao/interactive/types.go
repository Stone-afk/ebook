package interactive

import (
	"context"
	"gorm.io/gorm"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/dao/interactive/types.go -package=daomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/dao/mocks/interactive.mock.go
type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (Interactive, error)
	GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error)
	InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error)
}
