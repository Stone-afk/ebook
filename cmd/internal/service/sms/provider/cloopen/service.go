package cloopen

// 容联云短信的实现
// SDK文档:https://doc.yuntongxun.com/pe/5f029a06a80948a1006e7760

import (
	"context"
	"fmt"
	"github.com/cloopen/go-sms-sdk/cloopen"
	"log"
)

type Service struct {
	client *cloopen.SMS
	appId  string
}

func NewService(c *cloopen.SMS, addId string) *Service {
	return &Service{
		client: c,
		appId:  addId,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, data []string, numbers ...string) error {
	input := &cloopen.SendRequest{
		// 应用的APPID
		AppId: s.appId,
		// 模版ID
		TemplateId: tplId,
		// 模版变量内容 非必填
		Datas: data,
	}
	for _, number := range numbers {
		input.To = number

		resp, err := s.client.Send(input)
		if err != nil {
			return err
		}

		if resp.StatusCode != "000000" {
			log.Printf("response code: %s, msg: %s \n", resp.StatusCode, resp.StatusMsg)
			return fmt.Errorf("发送失败，code: %s, 原因：%s",
				resp.StatusCode, resp.StatusMsg)
		}
	}
	return nil
}
