package startup

import (
	"ebook/cmd/internal/handler"
	"ebook/cmd/internal/service/oauth2"
	"ebook/cmd/internal/service/oauth2/wechat"
	"ebook/cmd/pkg/logger"
	"net/http"
)

// InitPhantomWechatService 没啥用的虚拟的 wechatService
func InitPhantomWechatService(l logger.Logger) oauth2.Service {
	return wechat.NewService("ebook", "692jdHsogrsYqxaUK9fgxw", http.DefaultClient, l)
}

func NewWechatHandlerConfig() handler.WechatHandlerConfig {
	return handler.WechatHandlerConfig{
		Secure: false,
	}
}
