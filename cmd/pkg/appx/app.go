package appx

import (
	"ebook/cmd/pkg/ginx"
	"ebook/cmd/pkg/grpcx/server"
	"ebook/cmd/pkg/saramax"
	"ebook/cmd/pkg/schedulerx"
	"github.com/robfig/cron/v3"
)

// App 当在 wire 里面使用这个结构体的时候，要注意不是所有的服务都需要全部字段，
// 那么在 wire 的时候就不要使用 * 了
type App struct {
	Scheduler  schedulerx.Scheduler
	cron       *cron.Cron
	GRPCServer *server.Server
	WebServer  *ginx.Server
	Consumers  []saramax.Consumer
}
