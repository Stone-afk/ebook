package validator

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/migrator"
	"ebook/cmd/pkg/migrator/events"
	"github.com/ecodeclub/ekit/slice"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

// Validator T 必须实现了 Entity 接口
type Validator[T migrator.Entity] struct {
	baseValidator
	// 在这里加字段，比如说，在查询 base 根据什么列来排序，在 target 的时候，根据什么列来查询数据
	// 最极端的情况，是这样
	utime     int64
	batchSize int
	// 如果没有数据了，就睡眠
	// 如果不是正数，那么就说明直接返回，结束这一次的循环
	// 我很厌恶这种特殊值有特殊含义的做法，但是不得不搞
	sleepInterval time.Duration
	highLoad      *atomicx.Value[bool]
	fromBase      func(ctx context.Context, offset int) (T, error)
}

func NewValidator[T migrator.Entity](
	base *gorm.DB,
	target *gorm.DB,
	direction string,
	l logger.Logger,
	producer events.Producer) *Validator[T] {
	res := &Validator[T]{
		baseValidator: baseValidator{
			base:      base,
			target:    target,
			direction: direction,
			l:         l,
			producer:  producer,
		},
		highLoad: atomicx.NewValueOf[bool](false),
	}
	res.fromBase = res.fullFromBase
	return res
}

func (v *Validator[T]) Utime(utime int64) *Validator[T] {
	v.utime = utime
	return v
}

func (v *Validator[T]) SleepInterval(i time.Duration) *Validator[T] {
	v.sleepInterval = i
	return v
}

func (v *Validator[T]) SetFromBase(
	fromBase func(ctx context.Context, offset int) (T, error)) *Validator[T] {
	v.fromBase = fromBase
	return v
}

func (v *Validator[T]) Incr() *Validator[T] {
	v.fromBase = v.intrFromBase
	return v
}

func (v *Validator[T]) Validate(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(func() error {
		v.validateBaseToTarget(ctx)
		return nil
	})

	eg.Go(func() error {
		v.validateTargetToBase(ctx)
		return nil
	})
	return eg.Wait()
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
	offset := 0
	for {
		if v.highLoad.Load() {
			time.Sleep(10 * time.Minute)
		}
		src, err := v.fromBase(ctx, offset)
		switch err {
		case context.Canceled, context.DeadlineExceeded:
			// 超时或者被人取消了
			return
		case gorm.ErrRecordNotFound:
			// 比完了。没数据了，全量校验结束了
			// 同时支持全量校验和增量校验，你里就不能直接返回
			// 在这里要考虑：有些情况下，用户希望退出，有些情况下。用户希望继续
			// 当用户希望继续的时候要 sleep 一下
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
			continue
		case nil:
			// 查到了数据
			// 要去 target 里面找对应的数据
			var dst T
			err = v.target.WithContext(ctx).
				Where("id = ?", src.ID()).First(&dst).Error
			switch err {
			case context.Canceled, context.DeadlineExceeded:
				// 超时或者被人取消了
				return
			case gorm.ErrRecordNotFound:
				// 意味着，target 里面少了当前这条数据
				v.notify(ctx, src.ID(), events.InconsistentEventTypeTargetMissing)
			case nil:
				// 找到了。要开始比较
				// 怎么比较？
				// 能不能这么比？
				// 1. src == dst
				// 这是利用反射来比
				// 这个原则上是可以的。
				//if reflect.DeepEqual(src, dst) {
				//
				//}
				//var srcAny any = src
				//if c1, ok := srcAny.(interface {
				//	// 有没有自定义的比较逻辑
				//	CompareTo(c2 migrator.Entity) bool
				//}); ok {
				//	// 有，就用它的
				//	if !c1.CompareTo(dst) {
				//
				//	}
				//} else {
				//	// 没有，我就用反射
				//	if !reflect.DeepEqual(src, dst) {
				//
				//	}
				//}
				if !src.CompareTo(dst) {
					// 不相等, 这时候，上报给 Kafka，告知数据不一致
					v.notify(ctx, src.ID(), events.InconsistentEventTypeNEQ)
				}
			default:
				// 这里，要不要汇报，数据不一致？
				// 有两种做法：
				// 1. 认为，大概率数据是一致的，记录一下日志，下一条
				v.l.Error("查询 target 数据失败", logger.Error(err))
				// 2. 认为，出于保险起见，应该发送kafka通知报数据不一致，试着去修一下
				// 如果真的不一致了，执行修复
				// 如果假的不一致（也就是数据一致），也没事，就是多修了一次
				// 不好用哪个 InconsistentType
			}
		default:
			// 数据库错误
			v.l.Error("校验数据，查询 base 出错", logger.Error(err))
			time.Sleep(time.Second)
		}
		offset++
	}
}

// intrFromBase 增量校验
func (v *Validator[T]) intrFromBase(ctx context.Context, offset int) (T, error) {
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var src T
	// 找到了 base 中的数据
	// 例如 .Order("id DESC")，每次插入数据，就会导致你的 offset 不准了
	// 如果表没有 id 这个列怎么办？
	// 找一个类似的列，比如说 ctime (创建时间）
	err := v.base.WithContext(dbCtx).
		// 最好不要取等号
		Where("utime > ?", v.utime).
		Offset(offset).
		Order("utime ASC, id ASC").First(&src).Error
	// 按段取
	// WHERE utime >= ? LIMIT 10 ORDER BY UTIME
	// v.utime = srcList[len(srcList)].Utime()
	return src, err
}

// fullFromBase 全量校验
func (v *Validator[T]) fullFromBase(ctx context.Context, offset int) (T, error) {
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var src T
	// 找到了 base 中的数据
	// 例如 .Order("id DESC")，每次插入数据，就会导致你的 offset 不准了
	// 如果我的表没有 id 这个列怎么办？
	// 找一个类似的列，比如说 ctime (创建时间）
	err := v.base.WithContext(dbCtx).
		Offset(offset).
		Order("id").First(&src).Error
	return src, err
}

// 理论上来说，可以利用 count 来加速这个过程，
// 举个例子，假如说初始化目标表的数据是 昨天的 23:59:59 导出来的
// 那么可以 COUNT(*) WHERE ctime < 今天的零点，count 如果相等，就说明没删除
// 这一步大多数情况下效果很好，尤其是那些软删除的。
// 如果 count 不一致，那么接下来，理论上来说，还可以分段 count
// 比如说，先 count 第一个月的数据，一旦有数据删除了,你还得一条条查出来
// A utime=昨天 , A 在 base 里面，今天删了，A 在 target 里面，utime 还是昨天
// 这个地方，可以考虑不用 utime
// A 在删除之前，已经被修改过了，那么 A 在 target 里面的 utime 就是今天了
func (v *Validator[T]) validateTargetToBase(ctx context.Context) {
	// 先找 target，再找 base，找出 base 中已经被删除的
	// 理论上来说，就是 target 里面一条条找
	offset := 0
	for {
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		var dstTs []T
		err := v.target.WithContext(dbCtx).
			Where("utime > ?", v.utime).
			Select("id").
			Offset(offset).Limit(v.batchSize).
			Order("utime").First(dstTs).Error
		cancel()
		if len(dstTs) == 0 {
			// 没数据了。直接返回
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
			continue
		}
		switch err {
		case context.Canceled, context.DeadlineExceeded:
			// 超时或者被人取消了
			return
		case gorm.ErrRecordNotFound:
			// 没数据了。直接返回
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
			continue
		case nil:
			ids := slice.Map(dstTs, func(idx int, t T) int64 {
				return t.ID()
			})
			// 可以直接用 NOT IN
			var srcTs []T
			err = v.base.WithContext(ctx).
				Where("id IN ?", ids).Find(&srcTs).Error
			switch err {
			case context.Canceled, context.DeadlineExceeded:
				// 超时或者被人取消了
				return
			case gorm.ErrRecordNotFound:
				v.notifyBaseMissing(ctx, ids)
			case nil:
				srcIds := slice.Map(srcTs, func(idx int, t T) int64 {
					return t.ID()
				})
				// 计算差集
				// 也就是，src 里面的没有的
				diff := slice.DiffSet(ids, srcIds)
				v.notifyBaseMissing(ctx, diff)
			default:
				// 记录日志
				v.l.Error("查询 base失败", logger.Error(err))
			}
		default:
			// 记录日志，continue 掉
			v.l.Error("查询target 失败", logger.Error(err))
		}
		offset += v.batchSize
	}
}

func (v *Validator[T]) notifyBaseMissing(ctx context.Context, ids []int64) {
	for _, id := range ids {
		v.notify(ctx, id, events.InconsistentEventTypeBaseMissing)
	}
}

func (v *Validator[T]) notify(ctx context.Context, id int64, typ string) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	err := v.producer.ProduceInconsistentEvent(ctx, events.InconsistentEvent{
		ID:        id,
		Direction: v.direction,
		Type:      typ,
	})
	cancel()
	if err != nil {
		// 这又是一个问题
		// 怎么办？ 可以重试，但是重试也会失败，记日志，告警，手动去修
		// 直接忽略，下一轮修复和校验又会找出来
		v.l.Error("发送数据不一致的消息失败", logger.Error(err))
	}
}

type BatchValidator[T migrator.Entity] struct {
	baseValidator
	batchSize int
	utime     int64
	// 如果没有数据了，就睡眠
	// 如果不是正数，那么就说明直接返回，结束这一次的循环
	// 我很厌恶这种特殊值有特殊含义的做法，但是不得不搞
	sleepInterval time.Duration
}

func NewBatchValidator[T migrator.Entity](
	base *gorm.DB,
	target *gorm.DB,
	direction string,
	l logger.Logger,
	producer events.Producer,
) *BatchValidator[T] {
	return &BatchValidator[T]{
		baseValidator: baseValidator{
			base:      base,
			target:    target,
			direction: direction,
			l:         l,
			producer:  producer,
		},
		batchSize: 100,
		// 默认是全量校验，并且数据没了就结束
		sleepInterval: 0,
	}
}

func (v *BatchValidator[T]) Utime(utime int64) *BatchValidator[T] {
	v.utime = utime
	return v
}

func (v *BatchValidator[T]) SleepInterval(i time.Duration) *BatchValidator[T] {
	v.sleepInterval = i
	return v
}

// Validate 执行校验。
// 分成两步：
// 1. from => to
func (v *BatchValidator[T]) Validate(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(func() error {
		return v.baseToTarget(ctx)
	})
	eg.Go(func() error {
		return v.targetToBase(ctx)
	})
	return eg.Wait()
}

// baseToTarget 从 first 到 second 的验证
func (v *BatchValidator[T]) baseToTarget(ctx context.Context) error {
	offset := 0
	for {
		var src T
		// 这里假定主键的规范都是叫做 id，基本上大部分公司都有这种规范
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		err := v.base.WithContext(dbCtx).
			Order("id").
			Where("utime >= ?", v.utime).
			Offset(offset).First(&src).Error
		cancel()
		switch err {
		case gorm.ErrRecordNotFound:
			// 已经没有数据了
			if v.sleepInterval <= 0 {
				return nil
			}
			time.Sleep(v.sleepInterval)
			continue
		case context.Canceled, context.DeadlineExceeded:
			// 退出循环
			return nil
		case nil:
			v.dstDiff(ctx, src)
		default:
			v.l.Error("src => dst 查询源表失败", logger.Error(err))
		}
		offset++
	}
}

func (v *BatchValidator[T]) dstDiff(ctx context.Context, src T) {
	var dst T
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	err := v.target.WithContext(dbCtx).
		Where("id=?", src.ID()).First(&dst).Error
	cancel()
	// 这边要考虑不同的 error
	switch err {
	case gorm.ErrRecordNotFound:
		v.notify(src.ID(), events.InconsistentEventTypeTargetMissing)
	case nil:
		// 查询到了数据
		equal := src.CompareTo(dst)
		if !equal {
			v.notify(src.ID(), events.InconsistentEventTypeNEQ)
		}
	default:
		v.l.Error("src => dst 查询目标表失败", logger.Error(err))
	}
}

// targetToBase 反过来，执行 target 到 base 的验证
// 这是为了找出 dst 中多余的数据
func (v *BatchValidator[T]) targetToBase(ctx context.Context) error {
	// 这个我们只需要找出 src 中不存在的 id 就可以了
	offset := 0
	for {
		var ts []T
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		err := v.target.WithContext(dbCtx).Model(new(T)).Select("id").Offset(offset).
			Limit(v.batchSize).Find(&ts).Error
		cancel()
		switch err {
		case gorm.ErrRecordNotFound:
			if v.sleepInterval > 0 {
				time.Sleep(v.sleepInterval)
				// 在 sleep 的时候。不需要调整偏移量
				continue
			}
		case context.DeadlineExceeded, context.Canceled:
			return nil
		case nil:
			v.srcMissingRecords(ctx, ts)
		default:
			v.l.Error("dst => src 查询目标表失败", logger.Error(err))
		}
		if len(ts) < v.batchSize {
			// 数据没了
			return nil
		}
		offset += v.batchSize
	}
}

func (v *BatchValidator[T]) srcMissingRecords(ctx context.Context, ts []T) {
	ids := slice.Map(ts, func(idx int, src T) int64 {
		return src.ID()
	})
	dbCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	base := v.base.WithContext(dbCtx)
	var srcTs []T
	err := base.Select("id").Where("id IN ?", ids).Find(&srcTs).Error
	switch err {
	case gorm.ErrRecordNotFound:
		// 说明 ids 全部没有
		v.notifySrcMissing(ts)
	case nil:
		// 计算差集
		missing := slice.DiffSetFunc(ts, srcTs, func(src, dst T) bool {
			return src.ID() == dst.ID()
		})
		v.notifySrcMissing(missing)
	default:
		v.l.Error("dst => src 查询源表失败", logger.Error(err))
	}
}

func (v *BatchValidator[T]) notifySrcMissing(ts []T) {
	for _, t := range ts {
		v.notify(t.ID(), events.InconsistentEventTypeBaseMissing)
	}
}
