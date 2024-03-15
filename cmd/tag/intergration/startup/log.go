package startup

import (
	"ebook/cmd/pkg/logger"
	zaplogger "ebook/cmd/pkg/logger/zap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitLogger() logger.Logger {
	cfg := zap.NewDevelopmentConfig()
	err := viper.UnmarshalKey("log", &cfg)
	if err != nil {
		panic(err)
	}
	log, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return zaplogger.NewZapLogger(log)
}
