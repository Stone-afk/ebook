package startup

import (
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/logger/nop"
)

func InitLogger() logger.Logger {
	return nop.NewNopLogger()
}
