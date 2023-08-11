package middleware

import (
	"github.com/ecodeclub/ekit/set"
	"github.com/gin-gonic/gin"
)

type JWTLoginMiddlewareBuilder struct {
	publicPaths set.Set[string]
}

func NewJWTLoginMiddlewareBuilder() *JWTLoginMiddlewareBuilder {
	s := set.NewMapSet[string](3)
	s.Add("/users/signup")
	s.Add("/users/login_sms/code/send")
	s.Add("/users/login_sms")
	s.Add("/users/login")
	return &JWTLoginMiddlewareBuilder{
		publicPaths: s,
	}
}

func (l *JWTLoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要校验
		if l.publicPaths.Exist(ctx.Request.URL.Path) {
			return
		}
	}
}
