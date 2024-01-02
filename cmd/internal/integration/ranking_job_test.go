package integration

import (
	"ebook/cmd/internal/integration/startup"
	"ebook/cmd/internal/job"
	svcmocks "ebook/cmd/internal/service/mocks"
	"ebook/cmd/pkg/logger/nop"
	rlock "github.com/gotomicro/redis-lock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

// TestRankingJob 这个测试只是测试调度和 Redis 交互两个部分，但是不会真的测试计算的逻辑
func TestRankingJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	rdb := startup.InitRedis()
	svc := svcmocks.NewMockRankingService(ctrl)
	// 会调用三次
	svc.EXPECT().RankTopN(gomock.Any()).Times(3).Return(nil)
	j := job.NewRankingJob(svc, rlock.NewClient(rdb),
		nop.NewNopLogger(), time.Minute)
	c := cron.New(cron.WithSeconds())
	bd := job.NewCronJobBuilder(nop.NewNopLogger(),
		prometheus.SummaryOpts{
			Name: "test",
		})
	_, err := c.AddJob("@every 1s", bd.Build(j))
	require.NoError(t, err)
	c.Start()
	time.Sleep(time.Second * 3)
	ctx := c.Stop()
	<-ctx.Done()
}
