package failover

// 故障转移

import (
	"context"
	"ebook/cmd/sms/service"
	"errors"
	"sync/atomic"
)

type Service struct {
	services []service.Service
	idx      uint64
}

func NewService(services []service.Service) service.Service {
	return &Service{
		services: services,
	}
}

//func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
//	for _, svc := range s.services {
//		err := svc.Send(ctx, tpl, args, numbers...)
//		// 发送成功
//		if err == nil {
//			return nil
//		}
//		// 正常这边，输出日志
//		// 要做好监控
//		log.Println(err)
//	}
//	return errors.New("全部服务商都失败了")
//}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	// 取下一个节点来作为起始节点
	idx := atomic.AddUint64(&s.idx, 1)
	svcLength := uint64(len(s.services))
	length := idx + svcLength
	for i := idx; i < length; i++ {
		svc := s.services[int(i%svcLength)]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			// 超时/被上游Cancel
			return err
		default:
			// 输出日志
			// 这里可以对可能的报错继续扩展对应的分支
		}
	}
	return errors.New("全部服务商都失败了")
}
