package middleware

import (
	"github.com/ecodeclub/ekit/set"
	"github.com/gin-gonic/gin"
)

type JWTLoginMiddlewareBuilder struct {
	publicPaths set.Set[string]
}

func NewJWTLoginMiddlewareBuilder() *JWTLoginMiddlewareBuilder {
	return &JWTLoginMiddlewareBuilder{}
}

func (l *JWTLoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(context *gin.Context) {

	}
}
