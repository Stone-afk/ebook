package connpool

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"gorm.io/gorm"
)

const (
	patternDstOnly  = "DST_ONLY"
	patternSrcOnly  = "SRC_ONLY"
	patternDstFirst = "DST_FIRST"
	patternSrcFirst = "SRC_FIRST"
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
func (p *DoubleWritePool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
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

func (p *DoubleWritePool) BeginTx(ctx context.Context, opts *sql.TxOptions) (gorm.ConnPool, error) {
	pattern := p.pattern.Load()
	switch pattern {
	case patternSrcOnly:
		tx, err := p.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWritePoolTx{
			src:     tx,
			pattern: pattern,
		}, err
	case patternSrcFirst:
		srcTx, err := p.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		dstTx, err := p.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			// 记录日志，然后不做处理

			// 可以考虑回滚
			// err = srcTx.Rollback()
			// return err
		}
		return &DoubleWritePoolTx{
			src:     srcTx,
			dst:     dstTx,
			pattern: pattern,
		}, nil
	case patternDstOnly:
		tx, err := p.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		return &DoubleWritePoolTx{
			src:     tx,
			pattern: pattern,
		}, err
	case patternDstFirst:
		dstTx, err := p.dst.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}
		srcTx, err := p.src.(gorm.TxBeginner).BeginTx(ctx, opts)
		if err != nil {
			// 记录日志，然后不做处理

			// 可以考虑回滚
			// err = dstTx.Rollback()
			// return err
		}
		return &DoubleWritePoolTx{
			src:     srcTx,
			dst:     dstTx,
			pattern: pattern,
		}, nil
	default:
		return nil, errors.New("未知的双写模式")
	}
}

// ExecContext 在增量校验的时候，我能不能利用这个方法？
// 1.1 能不能从 query 里面抽取出来主键， WHERE id= xxx ，然后我就知道哪些数据被影响了？
// 1.2 可以尝试的思路是：用抽象语法树来分析 query， 而后找出 query 里面的条件，执行一个 SELECT，判定有哪些 id
// 1.2.1 UPDATE xx set b = xx WHERE a = 1
// 1.2.2 UPDATE xx set a = xx WHERE a = 1 LIMIT 10; DELETE from xxx WHERE aa OFFSET abc LIMIT cde
// 1.2.3 INSERT INTO ON CONFLICT, upsert 语句
func (p *DoubleWritePool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch p.pattern.Load() {
	case patternSrcOnly:
		return p.src.ExecContext(ctx, query, args...)
	case patternSrcFirst:
		res, err := p.src.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		if p.dst == nil {
			return res, err
		}
		_, err = p.dst.ExecContext(ctx, query, args...)
		if err != nil {
			// 记日志
			// dst 写失败，不被认为是失败
		}
		return res, err
	case patternDstOnly:
		return p.dst.ExecContext(ctx, query, args...)
	case patternDstFirst:
		res, err := p.dst.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		if p.src == nil {
			return res, err
		}
		_, err = p.src.ExecContext(ctx, query, args...)
		if err != nil {
			// 记日志
			// dst 写失败，不被认为是失败
		}
		return res, err
	default:
		// panic("未知的双写模式")
		return nil, errors.New("未知的双写模式")
	}
}

func (p *DoubleWritePool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch p.pattern.Load() {
	case patternSrcOnly, patternSrcFirst:
		return p.src.QueryContext(ctx, query, args...)
	case patternDstOnly, patternDstFirst:
		return p.dst.QueryContext(ctx, query, args...)
	default:
		// panic("未知的双写模式")
		return nil, errors.New("未知的双写模式")
	}
}

func (p *DoubleWritePool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch p.pattern.Load() {
	case patternSrcOnly, patternSrcFirst:
		return p.src.QueryRowContext(ctx, query, args...)
	case patternDstOnly, patternDstFirst:
		return p.dst.QueryRowContext(ctx, query, args...)
	default:
		// 这里有一个问题，怎么返回一个 error
		// unsafe 可以
		panic("未知的双写模式")
	}
}

func (p *DoubleWritePool) UpdatePattern(pattern string) {
	p.pattern.Store(pattern)
	// 能不能，有事务未提交的情况下，禁止修改
	// 能，但是性能问题比较严重，需要维持住一个已开事务的计数，要用锁了
}

type DoubleWritePoolTx struct {
	src     *sql.Tx
	dst     *sql.Tx
	pattern string
}

// Commit 和 PPT 不一致
func (p *DoubleWritePoolTx) Commit() error {
	switch p.pattern {
	case patternSrcOnly:
		return p.src.Commit()
	case patternSrcFirst:
		// 源库上的事务失败了，目标库要不要提交
		// commit 失败了怎么办？
		err := p.src.Commit()
		if err != nil {
			// 要不要提交？
			return err
		}
		if p.dst != nil {
			err = p.dst.Commit()
			if err != nil {
				// 记录日志
			}
		}
		return nil
	case patternDstOnly:
		return p.dst.Commit()
	case patternDstFirst:
		err := p.dst.Commit()
		if err != nil {
			// 要不要提交？
			return err
		}
		if p.src != nil {
			err = p.src.Commit()
			if err != nil {
				// 记录日志
			}
		}
		return nil
	default:
		return errUnknownPattern
	}
}

func (p *DoubleWritePoolTx) Rollback() error {
	switch p.pattern {
	case patternSrcOnly:
		return p.src.Rollback()
	case patternSrcFirst:
		// 源库上的事务失败了，目标库要不要提交
		// commit 失败嘞怎么办？
		err := p.src.Rollback()
		if err != nil {
			// 要不要提交？
			// 我个人觉得 可以尝试 rollback
			return err
		}
		if p.dst != nil {
			err = p.dst.Rollback()
			if err != nil {
				// 记录日志
			}
		}
		return nil
	case patternDstOnly:
		return p.dst.Rollback()
	case patternDstFirst:
		err := p.dst.Rollback()
		if err != nil {
			// 要不要提交？
			return err
		}
		if p.src != nil {
			err = p.src.Rollback()
			if err != nil {
				// 记录日志
			}
		}
		return nil
	default:
		return errUnknownPattern
	}
}

func (p *DoubleWritePoolTx) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	panic("implement me")
}

func (p *DoubleWritePoolTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	switch p.pattern {
	case patternSrcOnly:
		return p.src.ExecContext(ctx, query, args...)
	case patternSrcFirst:
		res, err := p.src.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		if p.dst == nil {
			return res, err
		}
		_, err = p.dst.ExecContext(ctx, query, args...)
		if err != nil {
			// 记日志
			// dst 写失败，不被认为是失败
		}
		return res, err
	case patternDstOnly:
		return p.dst.ExecContext(ctx, query, args...)
	case patternDstFirst:
		res, err := p.dst.ExecContext(ctx, query, args...)
		if err != nil {
			return res, err
		}
		if p.src == nil {
			return res, err
		}
		_, err = p.src.ExecContext(ctx, query, args...)
		if err != nil {
			// 记日志
			// dst 写失败，不被认为是失败
		}
		return res, err
	default:
		//panic("未知的双写模式")
		return nil, errors.New("未知的双写模式")
	}
}

func (p *DoubleWritePoolTx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	switch p.pattern {
	case patternSrcOnly, patternSrcFirst:
		return p.src.QueryContext(ctx, query, args...)
	case patternDstOnly, patternDstFirst:
		return p.dst.QueryContext(ctx, query, args...)
	default:
		// panic("未知的双写模式")
		return nil, errors.New("未知的双写模式")
	}
}

func (p *DoubleWritePoolTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	switch p.pattern {
	case patternSrcOnly, patternSrcFirst:
		return p.src.QueryRowContext(ctx, query, args...)
	case patternDstOnly, patternDstFirst:
		return p.dst.QueryRowContext(ctx, query, args...)
	default:
		panic("未知的双写模式")
		//return nil, errors.New("未知的双写模式")
	}
}
