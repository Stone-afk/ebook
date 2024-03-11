package jwt

import (
	"ebook/cmd/pkg/ginx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/bff/handler/jwt/types.go -package=jwtmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/bff/handler/jwt/mocks/jwt.mock.go
type Handler interface {
	SetLoginToken(ctx *gin.Context, uid int64) error
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	ClearToken(ctx *gin.Context) error
	CheckSession(ctx *gin.Context, ssid string) error
	ExtractToken(ctx *gin.Context) string
}

type RefreshClaims struct {
	UserId int64
	Ssid   string
	jwt.RegisteredClaims
}

type UserClaims = ginx.UserClaims