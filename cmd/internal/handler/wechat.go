package handler

import (
	"ebook/cmd/internal/service"
	"ebook/cmd/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
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

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	panic("")
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, state)
	// 要把我的 state 存好
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登录URL失败",
		})
		return
	}
	if err = h.setStateCookie(ctx, state); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	panic("")
}
