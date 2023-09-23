package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
}

type RefreshClaims struct {
	UserId int64
	Ssid   string
	jwt.RegisteredClaims
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明你自己的要放进去 token 里面的数据
	UserId    int64
	Ssid      string
	UserAgent string
}
