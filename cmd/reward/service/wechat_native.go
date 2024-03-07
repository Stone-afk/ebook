package service

import (
	pmtv1 "ebook/cmd/api/proto/gen/payment/v1"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/reward/repository"
)

type WechatNativeRewardService struct {
	client pmtv1.WechatPaymentServiceClient
	repo   repository.RewardRepository
	l      logger.Logger
}
