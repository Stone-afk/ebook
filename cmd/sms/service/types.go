package service

import "context"

// Service 发送短信的抽象
// 目前你可以理解为，这是一个为了适配不同的短信供应商的抽象
//
//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/sms/service/types.go -package=smsmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/sms/service/mocks/sms.mock.go
type Service interface {
	// Send tpl 短信模版， args 模版所需参数，numbers 电话号码和手机号码等
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
	//SendV1(ctx context.Context, tpl string, args []NamedArg, numbers ...string) error
	// 调用者需要知道实现者需要什么类型的参数，是 []string，还是 map[string]string
	//SendV2(ctx context.Context, tpl string, args any, numbers ...string) error
	//SendVV3(ctx context.Context, tpl string, args T, numbers ...string) error
}

type NamedArg struct {
	Val  string
	Name string
}
