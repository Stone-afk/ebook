package wechat

import (
	"ebook/cmd/payment/domain"
	"ebook/cmd/payment/events"
	"ebook/cmd/payment/repository"
	"ebook/cmd/pkg/logger"
	"errors"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
)

var errUnknownTransactionState = errors.New("未知的微信事务状态")

type NativePaymentService struct {
	svc       *native.NativeApiService
	appID     string
	mchID     string
	notifyURL string
	repo      repository.PaymentRepository
	l         logger.Logger

	// 在微信 native 里面，分别是
	// SUCCESS：支付成功
	// REFUND：转入退款
	// NOTPAY：未支付
	// CLOSED：已关闭
	// REVOKED：已撤销（付款码支付）
	// USERPAYING：用户支付中（付款码支付）
	// PAYERROR：支付失败(其他原因，如银行返回失败)
	nativeCBTypeToStatus map[string]domain.PaymentStatus
	producer             events.Producer
}

func NewNativePaymentService(svc *native.NativeApiService,
	repo repository.PaymentRepository,
	l logger.Logger,
	appid, mchid string) *NativePaymentService {
	return &NativePaymentService{
		l:     l,
		repo:  repo,
		svc:   svc,
		appID: appid,
		mchID: mchid,
		// 一般来说，这个都是固定的，基本不会变的
		// 这个从配置文件里面读取
		// 1. 测试环境 test.wechat.meoying.com
		// 2. 开发环境 dev.wecaht.meoying.com
		// 3. 线上环境 wechat.meoying.com
		// DNS 解析到腾讯云
		// wechat.tencent_cloud.meoying.com
		// DNS 解析到阿里云
		// wechat.ali_cloud.meoying.com
		notifyURL: "http://wechat.meoying.com/pay/callback",
		nativeCBTypeToStatus: map[string]domain.PaymentStatus{
			"SUCCESS":  domain.PaymentStatusSuccess,
			"PAYERROR": domain.PaymentStatusFailed,
			// 这个状态，有些人会考虑映射过去 PaymentStatusFailed
			"NOTPAY":     domain.PaymentStatusInit,
			"USERPAYING": domain.PaymentStatusInit,
			"CLOSED":     domain.PaymentStatusFailed,
			"REVOKED":    domain.PaymentStatusFailed,
			"REFUND":     domain.PaymentStatusRefund,
			// 其它状态你都可以加
		},
	}
}
