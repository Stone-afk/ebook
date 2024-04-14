package ioc

import (
	"ebook/cmd/im/service"
	"github.com/spf13/viper"
)

func GetSecret() service.Secret {
	type Config struct {
		secret string `yaml:"secret"`
	}
	var cfg Config
	err := viper.UnmarshalKey("OpenIM", &cfg)
	if err != nil {
		panic(err)
	}
	return service.Secret(cfg.secret)
}

func GetBaseHost() service.BaseHost {
	type Config struct {
		host string `yaml:"host"`
	}
	var cfg Config
	err := viper.UnmarshalKey("OpenIM", &cfg)
	if err != nil {
		panic(err)
	}
	return service.BaseHost(cfg.host)
}
