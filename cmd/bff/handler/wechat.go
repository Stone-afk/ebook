package handler

import (
	oauth2v1 "ebook/cmd/api/proto/gen/oauth2/v1"
	userv1 "ebook/cmd/api/proto/gen/user/v1"
	ijwt "ebook/cmd/bff/handler/jwt"
	"github.com/gin-gonic/gin"
)

var _ handler = (*OAuth2WechatHandler)(nil)

type OAuth2WechatHandler struct {
	// 这边也可以直接定义成 wechat.Service
	// 但是为了保持使用 mock 来测试，这里还是用了接口
	wechatSvc       oauth2v1.Oauth2ServiceClient
	userSvc         userv1.UserServiceClient
	stateCookieName string
	stateTokenKey   []byte
	ijwt.Handler
}

func (O OAuth2WechatHandler) RegisterRoutes(s *gin.Engine) {
	//TODO implement me
	panic("implement me")
}
