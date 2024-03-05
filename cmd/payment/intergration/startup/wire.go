//go:build wireinject

package startup

import (
	"ebook/cmd/payment/ioc"
	"ebook/cmd/payment/repository"
	"ebook/cmd/payment/repository/dao"
	"ebook/cmd/payment/service/wechat"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(InitLogger, InitTestDB, InitKafka)

var wechatNativeSvcSet = wire.NewSet(
	ioc.InitWechatClient,
	NewSyncProducer,
	dao.NewPaymentGORMDAO,
	repository.NewPaymentRepository,
	ioc.InitWechatNativeService,
	ioc.InitWechatConfig)

//go:generate wire
func InitWechatNativeService() *wechat.NativePaymentService {
	wire.Build(wechatNativeSvcSet, thirdPartySet)
	return new(wechat.NativePaymentService)
}
