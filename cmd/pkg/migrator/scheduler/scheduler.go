package scheduler

import (
	"ebook/cmd/pkg/gormx/connpool"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/migrator"
	"ebook/cmd/pkg/migrator/events"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sync"
)

// Scheduler 用来统一管理整个迁移过程
// 它不是必须的，可以理解为这是为了方便用户操作（和你理解）而引入的。
type Scheduler[T migrator.Entity] struct {
	l       logger.Logger
	lock    sync.Mutex
	src     *gorm.DB
	dst     *gorm.DB
	pool    *connpool.DoubleWritePool
	pattern string
	// 如果你要允许多个全量校验同时运行
	fulls map[string]func()
}

func NewScheduler[T migrator.Entity](
	l logger.Logger,
	src *gorm.DB,
	dst *gorm.DB,
	// 这个是业务用的 DoubleWritePool
	pool *connpool.DoubleWritePool,
	producer events.Producer) *Scheduler[T] {
	return &Scheduler[T]{
		l:       l,
		src:     src,
		dst:     dst,
		pattern: connpool.PatternSrcOnly,
		pool:    pool,
	}
}

// RegisterRoutes 这一个也不是必须的，就是你可以考虑利用配置中心，监听配置中心的变化
// 把全量校验，增量校验做成分布式任务，利用分布式任务调度平台来调度
func (s *Scheduler[T]) RegisterRoutes(server *gin.RouterGroup) {
	// 将这个暴露为 HTTP 接口
	// 可以配上对应的 UI
}

// ---- 下面是四个阶段 ---- //

type StartIncrRequest struct {
	Utime int64 `json:"utime"`
	// 毫秒数
	// json 不能正确处理 time.Duration 类型
	Interval int64 `json:"interval"`
}
