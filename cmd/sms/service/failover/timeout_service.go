package failover

import (
	"context"
	"ebook/cmd/sms/service"
	"sync/atomic"
)

type TimeoutService struct {
	// 你的服务商
	services []service.Service
	idx      int32
	// 连续超时的个数
	cnt int32

	// 阈值
	// 连续超时超过这个数字，就要切换
	threshold int32
}

func NewTimeoutService(services []service.Service, threshold int32) service.Service {
	return &TimeoutService{
		services:  services,
		threshold: threshold,
	}
}

func (s *TimeoutService) Send(ctx context.Context,
	tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&s.idx)
	cnt := atomic.LoadInt32(&s.cnt)
	if cnt > s.threshold {
		// 这里要切换，新的下标，往后挪了一个
		newIdx := (idx + 1) % int32(len(s.services))
		if atomic.CompareAndSwapInt32(&s.idx, idx, newIdx) {
			// 成功往后挪了一位
			atomic.StoreInt32(&s.cnt, 0)
		}
		// else 就是出现并发，别人换成功了
		//idx = newIdx
		idx = atomic.LoadInt32(&s.idx)
	}
	svc := s.services[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case nil:
		// 发送成功，连续失败的状态被打断了
		atomic.StoreInt32(&s.cnt, 0)
		return nil
	case context.DeadlineExceeded, context.Canceled:
		// 超时/被上游Cancel
		atomic.AddInt32(&s.cnt, 1)
		return err
	default:
		// 不知道什么错误
		// 你可以考虑，换下一个，语义则是：
		// - 超时错误，可能是偶发的，我尽量再试试
		// - 非超时，我直接下一个
		return err
	}
}
