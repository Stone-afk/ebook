//go:build wireinject
package startup

import (
	accountv1 "ebook/cmd/api/proto/gen/account/v1"
	pmtv1 "ebook/cmd/api/proto/gen/payment/v1"
	"ebook/cmd/reward/repository"
	"ebook/cmd/reward/repository/cache"
	"ebook/cmd/reward/repository/dao"
	"ebook/cmd/reward/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(InitTestDB, InitLogger, InitRedis)

//go:generate wire
func InitWechatNativeSvc(paymentClient pmtv1.WechatPaymentServiceClient,
	accountClient accountv1.AccountServiceClient) service.RewardService {
	wire.Build(
		thirdPartySet,
		dao.NewRewardGORMDAO,
		cache.NewRewardRedisCache,
		repository.NewRewardRepository,
		service.NewWechatNativeRewardService)
	return new(service.WechatNativeRewardService)
}
