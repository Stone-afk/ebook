package interactive

import (
	"context"
	"gorm.io/gorm"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/dao/interactive/types.go -package=daomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/dao/mocks/interactive.mock.go
type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
	DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
}
