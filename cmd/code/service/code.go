package service

import (
	"context"
	smsv1 "ebook/cmd/api/proto/gen/sms/v1"
	"ebook/cmd/code/repository"
	"fmt"
	"math/rand"
)

type codeService struct {
	repo   repository.CodeRepository
	smsSvc smsv1.SmsServiceClient
}

func NewCodeService(repo repository.CodeRepository, smsSvc smsv1.SmsServiceClient) CodeService {
	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// Send 发验证码，我需要什么参数？
func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	// 生成一个验证码
	code := svc.generateCode()
	// 塞进去 Redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		// 有问题
		return err
	}
	// 这前面成功了, 接下来发送出去
	// svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	//if err != nil {
	// 这个地方怎么办？
	// 这意味着，Redis 有这个验证码，但是不好意思，
	// 我能不能删掉这个验证码？
	// 你这个 err 可能是超时的 err，你都不知道，发出了没
	// 在这里重试
	// 要重试的话，初始化的时候，传入一个自己就会重试的 smsSvc
	//}
	//return svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	_, err = svc.smsSvc.Send(ctx, &smsv1.SmsSendRequest{
		TplId:   codeTplId,
		Args:    []string{code},
		Numbers: []string{phone},
	})
	return err
}

func (svc *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *codeService) generateCode() string {
	// 用随机数生成一个
	num := rand.Intn(999999)
	return fmt.Sprintf("%6d", num)
}
