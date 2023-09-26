package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	AccessTokenKey  = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	RefreshTokenKey = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvfx")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func (h *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	claims := ctx.MustGet("claims").(*UserClaims)
	return h.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid),
		"", time.Hour*24*7).Err()
}

func (h *RedisJWTHandler) SetLoginToken(ctx *gin.Context, userId int64) error {
	ssid := uuid.New().String()
	err := h.SetJWTToken(ctx, userId, ssid)
	if err != nil {
		return err
	}
	return h.setRefreshToken(ctx, userId, ssid)
}

func (h *RedisJWTHandler) SetJWTToken(ctx *gin.Context, userId int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserId:    userId,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AccessTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (h *RedisJWTHandler) setRefreshToken(ctx *gin.Context, userId int64, ssid string) error {
	claims := RefreshClaims{
		UserId: userId,
		Ssid:   ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RefreshTokenKey)
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}
