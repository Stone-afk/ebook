package dao

import (
	"context"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"gorm.io/gorm"
)

const (
	patternDstOnly  = "DST_ONLY"
	patternSrcOnly  = "SRC_ONLY"
	patternDstFirst = "DST_FIRST"
	patternSrcFirst = "SRC_FIRST"
)

type DoubleWriteDAO struct {
	src     InteractiveDAO
	dst     InteractiveDAO
	pattern *atomicx.Value[string]
}

func NewDoubleWriteDAO(src InteractiveDAO, dst InteractiveDAO) *DoubleWriteDAO {
	return &DoubleWriteDAO{
		src:     src,
		dst:     dst,
		pattern: atomicx.NewValueOf(patternSrcOnly),
	}
}

func NewDoubleWriteDAOV1(src *gorm.DB, dst *gorm.DB) *DoubleWriteDAO {
	return &DoubleWriteDAO{
		src:     NewGORMInteractiveDAO(src),
		dst:     NewGORMInteractiveDAO(dst),
		pattern: atomicx.NewValueOf(patternSrcOnly),
	}
}

func (d *DoubleWriteDAO) InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) Get(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) BatchIncrReadCnt(ctx context.Context, bizs []string, ids []int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) GetByIds(ctx context.Context, biz string, ids []int64) ([]Interactive, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleWriteDAO) UpdatePattern(pattern string) {
	d.pattern.Store(pattern)
}
