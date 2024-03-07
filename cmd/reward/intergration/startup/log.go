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
	l, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return zaplogger.NewZapLogger(l)
}
