package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/atomic"
	"io"
	"time"
)

const printLimit = 1024

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
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > printLimit {
			url = url[:1024]
		}
		accessLog := &AccessLog{
			Method: ctx.Request.Method,
			// URL 本身也可能很长
			Url: url,
		}
		if b.allowReqBody.Load() && ctx.Request.Body != nil {
			// Body 读完就没有了
			body, _ := ctx.GetRawData()
			// reader := io.NopCloser(bytes.NewBuffer(body))
			reader := io.NopCloser(bytes.NewReader(body))
			//ctx.Request.GetBody = func() (io.ReadCloser, error) {
			//	return reader, nil
			//}
			ctx.Request.Body = reader
			if len(body) > 1024 {
				body = body[:1024]
			}
			// 这其实是一个很消耗 CPU 和内存的操作
			// 因为会引起复制
			accessLog.ReqBody = string(body)
		}

		if b.allowRespBody.Load() {
			ctx.Writer = responseWriter{
				accessLog:      accessLog,
				ResponseWriter: ctx.Writer,
			}
		}

		defer func() {
			accessLog.Duration = time.Since(start).String()
			// accessLog.Duration = time.Now().Sub(start).String()
			b.loggerFunc(ctx, accessLog)
		}()

		// 执行到业务逻辑
		ctx.Next()
		//b.loggerFunc(ctx, al)
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

type responseWriter struct {
	accessLog *AccessLog
	gin.ResponseWriter
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.accessLog.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w responseWriter) Write(data []byte) (int, error) {
	w.accessLog.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriteString(data string) (int, error) {
	w.accessLog.RespBody = data
	return w.ResponseWriter.WriteString(data)
}
