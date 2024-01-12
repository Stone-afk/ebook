package connpool

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"gorm.io/gorm"
)

var errUnknownPattern = errors.New("未知的双写模式")

type DoubleWritePool struct {
	src     gorm.ConnPool
	dst     gorm.ConnPool
	pattern *atomicx.Value[string]
}

func NewDoubleWritePool(src gorm.ConnPool,
	dst gorm.ConnPool, pattern string) *DoubleWritePool {
	return &DoubleWritePool{
		src:     src,
		dst:     dst,
		pattern: atomicx.NewValueOf(pattern),
	}
}

// PrepareContext Prepare 的语句会进来这里
func (d *DoubleWritePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	// sql.Stmt 是一个结构体，没有办法说返回一个代表双写的 Stmt
	panic("implement me")
	//return nil, errors.New("双写模式下不支持")
	//switch d.pattern.Load() {
	//case PatternSrcOnly, PatternSrcFirst:
	//	return d.src.PrepareContext(ctx, query)
	//case PatternDstOnly, PatternDstFirst:
	//	return d.dst.PrepareContext(ctx, query)
	//default:
	//	panic("未知的双写模式")
	//	//return nil, errors.New("未知的双写模式")
	//}
}

func (d *DoubleWritePool) UpdatePattern(pattern string) {
	d.pattern.Store(pattern)
	// 能不能，有事务未提交的情况下，禁止修改
	// 能，但是性能问题比较严重，需要维持住一个已开事务的计数，要用锁了
}

type DoubleWritePoolTx struct {
	src     *sql.Tx
	dst     *sql.Tx
	pattern string
}
