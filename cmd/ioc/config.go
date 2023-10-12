package ioc

import (
	"context"
	"github.com/spf13/viper"
)

type Configer interface {
	GetString(ctx context.Context, key string) (string, error)
	MustGetString(ctx context.Context, key string) string
	GetStringOrDefault(ctc context.Context, key string, def string) string
	//Unmarshal()
}

type ViperConfigerAdapter struct {
	v *viper.Viper
}
