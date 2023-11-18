package middleware

import (
	ijwt "ebook/cmd/internal/handler/jwt"
	"github.com/ecodeclub/ekit/set"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type JWTLoginMiddlewareBuilder struct {
	publicPaths set.Set[string]
	ijwt.Handler
}

func NewJWTLoginMiddlewareBuilder(hdl ijwt.Handler) *JWTLoginMiddlewareBuilder {
	s := set.NewMapSet[string](3)
	s.Add("/users/signup")
	s.Add("/users/login_sms/code/send")
	s.Add("/users/login_sms")
	s.Add("/users/refresh_token")
	s.Add("/users/login")
	s.Add("/oauth2/wechat/authurl")
	s.Add("/oauth2/wechat/callback")
	s.Add("/test/metric")
	return &JWTLoginMiddlewareBuilder{
		publicPaths: s,
		Handler:     hdl,
	}
}

func (l *JWTLoginMiddlewareBuilder) IgnorePaths(path string) *JWTLoginMiddlewareBuilder {
	l.publicPaths.Add(path)
	return l
}

func (l *JWTLoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要校验
		if l.publicPaths.Exist(ctx.Request.URL.Path) {
			return
		}

		// Authorization 头部
		// 得到的格式 Bearer token
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			// 格式不对
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 如果是空字符串，你可以预期后面 Parse 就会报错
		tokenStr := l.ExtractToken(ctx)
		claims := ijwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			return ijwt.AccessTokenKey, nil
		})
		if err != nil || !token.Valid {
			// 不正确的 token
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		expireTime, err := claims.GetExpirationTime()
		if err != nil {
			// 拿不到过期时间
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if expireTime.Before(time.Now()) {
			// 已经过期
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if ctx.GetHeader("User-Agent") != claims.UserAgent {
			// 换了一个 User-Agent，可能是攻击者
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			// 系统错误或者用户已经主动退出登录了
			// 这里也可以考虑说，如果在 Redis 已经崩溃的时候，
			// 就不要去校验是不是已经主动退出登录了。
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 说明 token 是合法的
		// 我们把这个 token 里面的数据放到 ctx 里面，后面用的时候就不用再次 Parse 了
		ctx.Set("user", claims)
	}
}
