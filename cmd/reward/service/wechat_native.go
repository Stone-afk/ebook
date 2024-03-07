package service

import (
	"context"
	accountv1 "ebook/cmd/api/proto/gen/account/v1"
	pmtv1 "ebook/cmd/api/proto/gen/payment/v1"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/reward/domain"
	"ebook/cmd/reward/repository"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type WechatNativeRewardService struct {
	paymentSvc pmtv1.WechatPaymentServiceClient
	accountSvc accountv1.AccountServiceClient
	repo       repository.RewardRepository
	l          logger.Logger
}

func (s *WechatNativeRewardService) PreReward(ctx context.Context, r domain.Reward) (domain.CodeURL, error) {
	// 可以考虑缓存二维码，一旦我发现支付成功了，我就清除我的二维码
	cu, err := s.repo.GetCachedCodeURL(ctx, r)
	if err == nil {
		return cu, nil
	}
	r.Status = domain.RewardStatusInit
	rid, err := s.repo.CreateReward(ctx, r)
	if err != nil {
		return domain.CodeURL{}, err
	}
	resp, err := s.paymentSvc.NativePrePay(ctx, &pmtv1.PrePayRequest{
		Amt: &pmtv1.Amount{
			Total:    r.Amt,
			Currency: "CNY",
		},
		// 在这里记录分账信息
		//PayDetail: [
		//{"account": "platform", amt: 100,}
		//{"account": uid-123,, amt: 900}
		//],
		BizTradeNo:  s.bizTradeNO(rid),
		Description: fmt.Sprintf("打赏-%s", r.Target.BizName),
	})
	if err != nil {
		return domain.CodeURL{}, err
	}
	cu = domain.CodeURL{
		Rid: rid,
		URL: resp.CodeUrl,
	}
	// 当然可以异步了
	err1 := s.repo.CachedCodeURL(ctx, cu, r)
	if err1 != nil {
		// 记录日志
	}
	return cu, nil
}

func (s *WechatNativeRewardService) GetReward(ctx context.Context, rid, uid int64) (domain.Reward, error) {
	// 快路径
	r, err := s.repo.GetReward(ctx, rid)
	if err != nil {
		return domain.Reward{}, err
	}
	if r.Uid != uid {
		// 说明是非法查询
		return domain.Reward{}, errors.New("查询的打赏记录和打赏人对不上")
	}
	// 有可能，打赏记录，还是 Init 状态
	// 已经是完结状态
	if r.Completed() || ctx.Value("limited") == "true" {
		// 我已经知道你的支付结果了
		return r, nil
	}
	// 这个时候，考虑到支付到查询结果，搞一个慢路径
	// 有可能支付了，但是我 reward 本身没有收到通知
	// 直接查询 payment，
	// 只能解决，支付收到了，但是 reward 没收到
	// 降级状态，限流状态，熔断状态，不要走慢路径
	resp, err := s.paymentSvc.GetPayment(ctx, &pmtv1.GetPaymentRequest{
		BizTradeNo: s.bizTradeNO(rid),
	})
	if err != nil {
		// 这边我们直接返回从数据库查询的数据
		s.l.Error("慢路径查询支付结果失败",
			logger.Int64("rid", r.Id), logger.Error(err))
		return r, nil
	}
	// 更新状态
	switch resp.Status {
	case pmtv1.PaymentStatus_PaymentStatusFailed:
		r.Status = domain.RewardStatusFailed
	case pmtv1.PaymentStatus_PaymentStatusInit:
		r.Status = domain.RewardStatusInit
	case pmtv1.PaymentStatus_PaymentStatusSuccess:
		r.Status = domain.RewardStatusPayed
	case pmtv1.PaymentStatus_PaymentStatusRefund:
		// 理论上来说不可能出现这个，直接设置为失败
		r.Status = domain.RewardStatusFailed
	}
	err = s.repo.UpdateStatus(ctx, rid, r.Status)
	if err != nil {
		s.l.Error("更新本地打赏状态失败",
			logger.Int64("rid", r.Id), logger.Error(err))
		return r, nil
	}
	return r, nil
}

func (s *WechatNativeRewardService) UpdateReward(ctx context.Context, bizTradeNO string, status domain.RewardStatus) error {
	rid := s.toRid(bizTradeNO)
	err := s.repo.UpdateStatus(ctx, rid, status)
	if err != nil {
		return err
	}
	// 完成了支付，准备入账
	// 已经支付成功了
	if status == domain.RewardStatusPayed {
		r, er := s.repo.GetReward(ctx, rid)
		if er != nil {
			return er
		}
		// ebook 抽成
		// 0.1 可以写到数据库里面，因为不同的商家可能有不同的议价权
		// 订单计算总价 + 分账
		// 分账要小心精度问题
		eAmt := int64(float64(r.Amt) * 0.1)
		_, err = s.accountSvc.Credit(ctx, &accountv1.CreditRequest{
			Biz:   "reward",
			BizId: rid,
			Items: []*accountv1.CreditItem{
				{
					AccountType: accountv1.AccountType_AccountTypeReward,
					// 虽然可能为 0，但是也要记录出来
					Amt:      eAmt,
					Currency: "CNY",
				},
				{
					//  严格来说 Account 需要从 user / 账户系统获取，默认为 r.Uid
					Account:     r.Uid,
					Uid:         r.Uid,
					AccountType: accountv1.AccountType_AccountTypeReward,
					Amt:         r.Amt - eAmt,
					Currency:    "CNY",
				},
			},
		})
		if err != nil {
			s.l.Error("入账失败了，快来修数据啊！！！",
				logger.String("biz_trade_no", bizTradeNO),
				logger.Error(err))
			// 做好监控和告警，这里
			// 引入自动修复功能
			// 如果没有 24小时值班 + 自动修复 + 异地容灾备份（随机演练）
			// 然后面试官又吹牛逼说自己的可用性有 9999，你就可以认为，他在扯淡。
			return err
		}
	}
	return nil
}

func (s *WechatNativeRewardService) bizTradeNO(rid int64) string {
	return fmt.Sprintf("reward-%d", rid)
}

func (s *WechatNativeRewardService) toRid(tradeNO string) int64 {
	ridStr := strings.Split(tradeNO, "-")
	val, _ := strconv.ParseInt(ridStr[1], 10, 64)
	return val
}

func NewWechatNativeRewardService(
	paymentSvc pmtv1.WechatPaymentServiceClient,
	repo repository.RewardRepository,
	l logger.Logger,
	accountSvc accountv1.AccountServiceClient,
) RewardService {
	return &WechatNativeRewardService{paymentSvc: paymentSvc, repo: repo, l: l, accountSvc: accountSvc}
}
