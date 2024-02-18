package circuitbreaker

import (
	"context"
	"github.com/go-kratos/aegis/circuitbreaker"
	"google.golang.org/grpc"
)

type InterceptorBuilder struct {
	breaker circuitbreaker.CircuitBreaker
	// 设置标记位
	// 假如说我们考虑使用随机数 + 阈值的回复方式
	// 触发熔断的时候，直接将 threshold 置为0
	// 后续等一段时间，将 theshold 调整为 1，判定请求有没有问题
	threshold int
}

func (b *InterceptorBuilder) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		panic("")
	}
}

func (b *InterceptorBuilder) allow() bool {
	// 这边就套用我们之前在短信里面讲的，判定节点是否健康的各种做法
	// 从prometheus 里面拿数据判定
	// prometheus.DefaultGatherer.Gather()
	return false
}
