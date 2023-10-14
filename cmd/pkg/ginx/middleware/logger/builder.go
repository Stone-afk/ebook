package logger

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
)

// MiddlewareBuilder 注意点：
// 1. 小心日志内容过多。URL 可能很长，请求体，响应体都可能很大，你要考虑是不是完全输出到日志里面
// 2. 考虑 1 的问题，以及用户可能换用不同的日志框架，所以要有足够的灵活性
// 3. 考虑动态开关，结合监听配置文件，要小心并发安全
type MiddlewareBuilder struct {
	// allowRespBody bool 并发不安全，需要加锁，然而通常加锁对性能损耗大，所以考虑原子操作
	allowReqBody  *atomic.Bool
	allowRespBody *atomic.Bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc:    fn,
		allowReqBody:  atomic.NewBool(false),
		allowRespBody: atomic.NewBool(false),
	}
}

func (b *MiddlewareBuilder) AllowReqBody(ok bool) *MiddlewareBuilder {
	b.allowReqBody.Store(ok)
	return b
}

func (b *MiddlewareBuilder) AllowRespBody(ok bool) *MiddlewareBuilder {
	b.allowRespBody.Store(ok)
	return b
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

type AccessLog struct {
	// HTTP 请求的方法
	Method string
	// Url 整个请求 URL
	Url      string
	Duration string
	ReqBody  string
	RespBody string
	Status   int
}
