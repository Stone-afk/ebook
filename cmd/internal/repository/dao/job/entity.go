package job

const (
	// 等待被调度，意思就是没有人正在调度
	jobStatusWaiting = iota
	// 已经被 goroutine 抢占了
	jobStatusRunning
	// 不再需要调度了，比如说被终止了，或者被删除了。
	// 我们这里没有严格区分这两种情况的必要性
	// 暂停调度
	jobStatusPaused
)

type Job struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Name     string `gorm:"type:varchar(256);unique"`
	Cfg      string
	Version  int64
	Executor string
	// cron 表达式
	Expression string
	// 另外一个问题，定时任务，我怎么知道，已经到时间了呢？
	// NextTime 下一次被调度的时间
	// next_time <= now 这样一个查询条件
	// and status = 0
	// 要建立索引
	// 更加好的应该是 next_time 和 status 的联合索引
	NextTime int64 `gorm:"index"`
	// 第一个问题：哪些任务可以抢？哪些任务已经被人占着？哪些任务永远不会被运行
	// 用状态来标记
	Status int
	Ctime  int64
	Utime  int64
}
