package startup

import (
	"ebook/cmd/internal/handler"
	"ebook/cmd/internal/service/oauth2"
	"ebook/cmd/internal/service/oauth2/wechat"
	"ebook/cmd/pkg/logger"
	"net/http"
	"os"
)

func InitWechatService(l logger.Logger) oauth2.Service {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		//panic("没有找到环境变量 WECHAT_APP_ID ")
		appId = "ebook"
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		// panic("没有找到环境变量 WECHAT_APP_SECRET")
		appKey = "692jdHsogrsYqxaUK9fgxw"
	}
	return wechat.NewService(appId, appKey, http.DefaultClient, l)
}

func NewWechatHandlerConfig() handler.WechatHandlerConfig {
	return handler.WechatHandlerConfig{
		Secure: false,
	}
}
