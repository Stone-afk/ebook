package ioc

import (
	wechat2 "ebook/cmd/oauth2/service/wechat"
	"ebook/cmd/pkg/logger"
	"github.com/spf13/viper"
)

func InitPrometheus(log logger.Logger) wechat2.Service {
	svc := InitService(log)
	type Config struct {
		NameSpace  string `yaml:"nameSpace"`
		Subsystem  string `yaml:"subsystem"`
		InstanceID string `yaml:"instanceId"`
		Name       string `yaml:"name"`
	}
	var cfg Config
	err := viper.UnmarshalKey("prometheus", &cfg)
	if err != nil {
		panic(err)
	}
	return wechat2.NewPrometheusDecorator(svc, cfg.NameSpace, cfg.Subsystem, cfg.InstanceID, cfg.Name)
}
