package logging

import (
	"ebook/cmd/pkg/grpcx/interceptors"
	"ebook/cmd/pkg/logger"
)

type InterceptorBuilder struct {
	// 如果你要非常通用
	l logger.Logger
	//fn func(msg string, fields...logger.Field)
	interceptors.Builder

	reqBody  bool
	respBody bool
}
