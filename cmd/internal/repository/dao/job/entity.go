package job

const (
	// 等待被调度，意思就是没有人正在调度
	jobStatusWaiting = iota
	// 已经被 goroutine 抢占了
	jobStatusRunning
	// 不再需要调度了，比如说被终止了，或者被删除了。
	// 我们这里没有严格区分这两种情况的必要性
	jobStatusEnd
)

type Job struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Name     string `gorm:"type:varchar(256);unique"`
	Cfg      string
	Version  int64
	NextTime int64 `gorm:"index"`
	Status   int
	Ctime    int64
	Utime    int64
}
