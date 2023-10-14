package zap

import (
	"ebook/cmd/pkg/logger"
	"go.uber.org/zap"
)

type zapLogger struct {
	log *zap.Logger
}

func NewZapLogger(log *zap.Logger) logger.Logger {
	return &zapLogger{
		log: log,
	}
}

func (z *zapLogger) toZapFields(args []logger.Field) []zap.Field {
	res := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		res = append(res, zap.Any(arg.Key, arg.Value))
	}
	return res
}

func (z *zapLogger) Debug(msg string, args ...logger.Field) {
	z.log.Debug(msg, z.toZapFields(args)...)
}

func (z *zapLogger) Info(msg string, args ...logger.Field) {
	z.log.Info(msg, z.toZapFields(args)...)
}

func (z *zapLogger) Warn(msg string, args ...logger.Field) {
	z.log.Warn(msg, z.toZapFields(args)...)
}

func (z *zapLogger) Error(msg string, args ...logger.Field) {
	z.log.Error(msg, z.toZapFields(args)...)
}
