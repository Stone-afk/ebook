package zap

import (
	"ebook/cmd/pkg/logger"
	"go.uber.org/zap"
)

type zapLogger struct {
	l *zap.Logger
}

func NewZapLogger(l *zap.Logger) logger.Logger {
	return &zapLogger{
		l: l,
	}
}

func (z *zapLogger) Debug(msg string, args ...logger.Field) {
	panic("")
}

func (z *zapLogger) Info(msg string, args ...logger.Field) {
	panic("")
}

func (z *zapLogger) Warn(msg string, args ...logger.Field) {
	panic("")
}

func (z *zapLogger) Error(msg string, args ...logger.Field) {
	panic("")
}
