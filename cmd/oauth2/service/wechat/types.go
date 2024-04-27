package wechat

import (
	"context"
	"ebook/cmd/oauth2/domain/wechat"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/oauth2/service/wechat/types.go -package=wechatmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/oauth2/service/wechat/mocks/svc.mock.go
type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	// VerifyCode 目前大部分公司的 OAuth2 平台都差不多的设计
	// 返回一个 unionId。这个你可以理解为，在第三方平台上的 unionId
	// 你也可以考虑使用 openId 来替换
	// 一家公司如果有很多应用，不同应用都有自建的用户系统
	// 那么 openId 可能更加合适
	VerifyCode(ctx context.Context, code string) (wechat.WechatInfo, error)
}