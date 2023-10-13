package nop

import "ebook/cmd/pkg/logger"

type nopLogger struct {
}

func NewNopLogger() logger.Logger {
	return &nopLogger{}
}

func (n *nopLogger) Debug(msg string, args ...logger.Field) {
}

func (n *nopLogger) Info(msg string, args ...logger.Field) {
}

func (n *nopLogger) Warn(msg string, args ...logger.Field) {
}

func (n *nopLogger) Error(msg string, args ...logger.Field) {
}
