package ginx

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	Id        int64
	UserAgent string
	Ssid      string
	jwt.RegisteredClaims
}
