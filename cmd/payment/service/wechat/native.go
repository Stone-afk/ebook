package wechat

import (
	"context"
	"ebook/cmd/payment/domain"
	"ebook/cmd/payment/events"
	"ebook/cmd/payment/repository"
	"ebook/cmd/pkg/logger"
	"errors"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"time"
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
	producer events.Producer,
	l logger.Logger,
	appid, mchid string) *NativePaymentService {
	return &NativePaymentService{
		l:        l,
		repo:     repo,
		svc:      svc,
		appID:    appid,
		mchID:    mchid,
		producer: producer,
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

// Prepay 为了拿到扫码支付的二维码
func (s *NativePaymentService) Prepay(ctx context.Context, pmt domain.Payment) (string, error) {
	// 唯一索引冲突
	// 业务方唤起了支付，但是没付，下一次再过来，应该换 BizTradeNO
	err := s.repo.AddPayment(ctx, pmt)
	if err != nil {
		return "", err
	}
	//sn := uuid.New().String()
	resp, result, err := s.svc.Prepay(ctx, native.PrepayRequest{
		Appid:       core.String(s.appID),
		Mchid:       core.String(s.mchID),
		Description: core.String(pmt.Description),
		// 这个地方是有讲究的
		// 选择1：业务方直接给我，我透传，我啥也不干
		// 选择2：业务方给我它的业务标识，我自己生成一个 - 担忧出现重复
		// 注意，不管你是选择 1 还是选择 2，业务方都一定要传给你（webook payment）一个唯一标识
		// Biz + BizTradeNo 唯一， biz + biz_id
		OutTradeNo: core.String(pmt.BizTradeNO),
		NotifyUrl:  core.String(s.notifyURL),
		// 设置三十分钟有效
		TimeExpire: core.Time(time.Now().Add(time.Minute * 30)),
		Amount: &native.Amount{
			Total:    core.Int64(pmt.Amt.Total),
			Currency: core.String(pmt.Amt.Currency),
		},
	})
	s.l.Debug("微信prepay响应",
		logger.Field{Key: "result", Value: result},
		logger.Field{Key: "resp", Value: resp})
	if err != nil {
		return "", err
	}
	// 这里可以考虑引入另外一个状态，也就是代表你已经调用了第三方支付，正在等回调的状态
	// 但是这个状态意义不是很大。
	// 因为在考虑兜底（定时比较数据）的时候，不管有没有调用第三方支付，
	// 都要问一下第三方支付这个
	return *resp.CodeUrl, nil
}

func (s *NativePaymentService) SyncWechatInfo(ctx context.Context, bizTradeNO string) error {
	txn, _, err := s.svc.QueryOrderByOutTradeNo(ctx, native.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(bizTradeNO),
		Mchid:      core.String(s.mchID),
	})
	if err != nil {
		return err
	}
	return s.updateByTxn(ctx, txn)
}

func (s *NativePaymentService) GetPayment(ctx context.Context, bizTradeId string) (domain.Payment, error) {
	// 在这里，能不能设计一个慢路径？如果要是不知道支付结果，就去微信里面查一下？
	// 或者异步查一下？
	return s.repo.GetPayment(ctx, bizTradeId)
}

func (s *NativePaymentService) FindExpiredPayment(ctx context.Context, offset, limit int, t time.Time) ([]domain.Payment, error) {
	return s.repo.FindExpiredPayment(ctx, offset, limit, t)
}

func (s *NativePaymentService) HandleCallback(ctx context.Context, txn *payments.Transaction) error {
	return s.updateByTxn(ctx, txn)
}

func (s *NativePaymentService) updateByTxn(ctx context.Context, txn *payments.Transaction) error {
	// 搞一个 status 映射的 map
	status, ok := s.nativeCBTypeToStatus[*txn.TradeType]
	if !ok {
		return fmt.Errorf("%w, %s", errUnknownTransactionState, *txn.TradeState)
	}
	pmt := domain.Payment{
		BizTradeNO: *txn.OutTradeNo,
		TxnID:      *txn.TransactionId,
		Status:     status,
	}
	// 核心就是更新数据库状态
	err := s.repo.UpdatePayment(ctx, pmt)
	if err != nil {
		// 这里有一个小问题，就是如果超时了的话，都不知道更新成功了没
		return err
	}
	// 就是处于结束状态
	// 发送消息，有结果了总要通知业务方
	// 这里有很多问题，核心就是部分失败问题，其次还有重复发送问题
	err1 := s.producer.ProducePaymentEvent(ctx, events.PaymentEvent{
		BizTradeNO: *txn.OutTradeNo,
		Status:     status.AsUint8(),
	})
	if err1 != nil {
		// 加监控加告警，立刻手动修复，或者自动补发
		s.l.Error("发送支付事件失败", logger.Error(err),
			logger.String("biz_trade_no", pmt.BizTradeNO))
	}
	// 虽然发送事件失败，但是数据库记录了，所以可以返回 Nil
	return nil
}
