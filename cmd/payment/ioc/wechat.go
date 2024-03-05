package ioc

import (
	"context"
	"ebook/cmd/payment/events"
	"ebook/cmd/payment/repository"
	"ebook/cmd/payment/service/wechat"
	"ebook/cmd/pkg/logger"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"os"
)

func InitWechatClient(cfg WechatConfig) *core.Client {
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(
		// 注意这个文件我没有上传，所以你需要准备一个
		cfg.KeyPath,
	)
	if err != nil {
		panic(err)
	}
	// 使用商户私钥等初始化 client
	client, err := core.NewClient(context.Background(), option.WithWechatPayAutoAuthCipher(
		cfg.MchID, cfg.MchSerialNum, mchPrivateKey, cfg.MchKey))
	if err != nil {
		panic(err)
	}
	return client
}

func InitWechatNativeService(
	cli *core.Client,
	repo repository.PaymentRepository,
	l logger.Logger,
	producer events.Producer,
	cfg WechatConfig) *wechat.NativePaymentService {
	return wechat.NewNativePaymentService(&native.NativeApiService{
		Client: cli,
	}, repo, producer, l, cfg.AppID, cfg.MchID)
}

func InitWechatConfig() WechatConfig {
	return WechatConfig{
		AppID:        os.Getenv("WEPAY_APP_ID"),
		MchID:        os.Getenv("WEPAY_MCH_ID"),
		MchKey:       os.Getenv("WEPAY_MCH_KEY"),
		MchSerialNum: os.Getenv("WEPAY_MCH_SERIAL_NUM"),
		CertPath:     "./config/cert/apiclient_cert.pem",
		KeyPath:      "./config/cert/apiclient_key.pem",
	}
}

type WechatConfig struct {
	AppID        string
	MchID        string
	MchKey       string
	MchSerialNum string

	// 证书
	CertPath string
	KeyPath  string
}
