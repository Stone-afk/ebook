package startup

import (
	"ebook/cmd/pkg/logger"
	zaplogger "ebook/cmd/pkg/logger/zap"
	"go.uber.org/zap"
)

func InitLogger() logger.Logger {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return zaplogger.NewZapLogger(log)
}
