package test

import (
	feedv1 "ebook/cmd/api/proto/gen/feed/v1"
	followMocks "ebook/cmd/api/proto/gen/followrelation/v1/mocks"
	"ebook/cmd/feed/grpc"
	"ebook/cmd/feed/ioc"
	"ebook/cmd/feed/repository"
	"ebook/cmd/feed/repository/cache"
	"ebook/cmd/feed/repository/dao"
	"ebook/cmd/feed/service"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
)

func InitGrpcServer(t *testing.T) (feedv1.FeedSvcServer, *followMocks.MockFollowServiceClient, *gorm.DB) {
	loggerV1 := ioc.InitLogger()
	db := ioc.InitDB(loggerV1)
	feedPullEventDAO := dao.NewFeedPullEventDAO(db)
	feedPushEventDAO := dao.NewFeedPushEventDAO(db)
	cmdable := ioc.InitRedis()
	feedEventCache := cache.NewFeedEventCache(cmdable)
	feedEventRepo := repository.NewFeedEventRepo(feedPullEventDAO, feedPushEventDAO, feedEventCache)
	mockCtrl := gomock.NewController(t)
	followClient := followMocks.NewMockFollowServiceClient(mockCtrl)
	v := ioc.RegisterHandler(feedEventRepo, followClient)
	feedService := service.NewFeedService(feedEventRepo, v)
	feedEventGrpcSvc := grpc.NewFeedEventServiceServer(feedService)
	return feedEventGrpcSvc, followClient, db
}
