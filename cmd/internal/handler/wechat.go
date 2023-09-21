package handler

import (
	"ebook/cmd/internal/service"
	"ebook/cmd/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
)

type WechatHandlerConfig struct {
	Secure bool
	//StateKey
}

type OAuth2WechatHandler struct {
	svc     wechat.Service
	userSvc service.UserService

	stateKey []byte
	cfg      WechatHandlerConfig
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	panic("")
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	panic("")
}
