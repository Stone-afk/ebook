package handler

import (
	oauth2v1 "ebook/cmd/api/proto/gen/oauth2/v1"
	userv1 "ebook/cmd/api/proto/gen/user/v1"
	ijwt "ebook/cmd/bff/handler/jwt"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
)

var (
	_               handler = (*OAuth2WechatHandler)(nil)
	stateCookieName         = "jwt-state"
)

type WechatHandlerConfig struct {
	Secure bool
	//StateKey
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}

type OAuth2WechatHandler struct {
	// 这边也可以直接定义成 wechat.Service
	// 但是为了保持使用 mock 来测试，这里还是用了接口
	wechatSvc       oauth2v1.Oauth2ServiceClient
	userSvc         userv1.UserServiceClient
	cfg             WechatHandlerConfig
	stateCookieName string
	stateTokenKey   []byte
	ijwt.Handler
}

func (h *OAuth2WechatHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	// 这边用 Any 万无一失
	g.Any("/callback", h.Callback)
}

func NewOAuth2WechatHandler(svc oauth2v1.Oauth2ServiceClient,
	userSvc userv1.UserServiceClient,
	jwtHdl ijwt.Handler,
	cfg WechatHandlerConfig) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		wechatSvc: svc,
		userSvc:   userSvc,
		// 万一后续我们要改，也可以做成可配置的。
		stateCookieName: "jwt-state",
		stateTokenKey:   []byte("95osj3fUD7foxmlYdDbncXz4VD2igvf1"),
		cfg:             cfg,
		Handler:         jwtHdl,
	}
}

// setStateCookie 只有微信这里用，所以定义在这里
func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
	})
	tokenStr, err := token.SignedString(h.stateTokenKey)
	if err != nil {
		return err
	}
	ctx.SetCookie(stateCookieName, tokenStr,
		600,
		// 限制在只能在这里生效。
		"/oauth2/wechat/callback",
		// 这边把 HTTPS 协议禁止了。不过在生产环境中要开启。
		"", h.cfg.Secure, true)
	return nil
}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.wechatSvc.AuthURL(ctx, &oauth2v1.AuthURLRequest{
		State: state,
	})
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
	code := ctx.Query("code")
	err := h.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登录失败",
		})
		return
	}
	info, err := h.wechatSvc.VerifyCode(ctx, &oauth2v1.VerifyCodeRequest{
		Code: code,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 这里怎么办？
	// 从 userService 里面拿 uid
	u, err := h.userSvc.FindOrCreateByWechat(ctx,
		&userv1.FindOrCreateByWechatRequest{
			Info: &userv1.WechatInfo{
				OpenId:  info.OpenId,
				UnionId: info.UnionId,
			},
		})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	err = h.SetLoginToken(ctx, u.User.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
	// 验证微信的 code
}

func (h *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	// 校验一下我的 state
	tokenStr, err := ctx.Cookie(stateCookieName)
	if err != nil {
		return fmt.Errorf("拿不到 state 的 cookie, %w", err)
	}
	var sc StateClaims
	token, err := jwt.ParseWithClaims(tokenStr, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateTokenKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("token 已经过期了, %w", err)
	}
	if sc.State != state {
		return errors.New("state 不相等")
	}
	return nil
}
