package cache

import (
	"context"
	_ "embed"
)

var (
	//go:embed lua/interative_incr_cnt.lua
	luaIncrCnt string
)

const (
	fieldReadCnt    = "read_cnt"
	fieldCollectCnt = "collect_cnt"
	fieldLikeCnt    = "like_cnt"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/cache/interactive.go -package=cachemocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/cache/mocks/interactive.mock.go
type InteractiveCache interface {
	// IncrReadCntIfPresent 如果在缓存中有对应的数据，就 +1
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
}
