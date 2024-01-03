package validator

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/migrator"
	"gorm.io/gorm"
)

// Validator T 必须实现了 Entity 接口
type Validator[T migrator.Entity] struct {
	// 校验，以 XXX 为准，
	base *gorm.DB
	// 校验的是谁的数据
	target    *gorm.DB
	l         logger.Logger
	direction string

	batchSize int
}

// <utime, id> 然后执行 SELECT * FROM xx WHERE utime > ? ORDER BY id
// 索引排序，还是内存排序？

// Validate 调用者可以通过 ctx 来控制校验程序退出
// 全量校验，是不是一条条比对？
// 所以要从数据库里面一条条查询出来
// utime 上面至少要有一个索引，并且 utime 必须是第一列
// <utime, col1, col2>, <utime> 这种可以
// <col1, utime> 这种就不可以
func (v *Validator[T]) validateBaseToTarget(ctx context.Context) {
	panic("")
}

// 理论上来说，可以利用 count 来加速这个过程，
// 我举个例子，假如说你初始化目标表的数据是 昨天的 23:59:59 导出来的
// 那么你可以 COUNT(*) WHERE ctime < 今天的零点，count 如果相等，就说明没删除
// 这一步大多数情况下效果很好，尤其是那些软删除的。
// 如果 count 不一致，那么接下来，你理论上来说，还可以分段 count
// 比如说，我先 count 第一个月的数据，一旦有数据删除了，你还得一条条查出来
// A utime=昨天
// A 在 base 里面，今天删了，A 在 target 里面，utime 还是昨天
// 这个地方，可以考虑不用 utime
// A 在删除之前，已经被修改过了，那么 A 在 target 里面的 utime 就是今天了
func (v *Validator[T]) validateTargetToBase(ctx context.Context) {
	panic("")
}
