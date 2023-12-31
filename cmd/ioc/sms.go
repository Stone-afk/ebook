package ioc

import (
	"ebook/cmd/internal/service/sms"
	"ebook/cmd/internal/service/sms/memory"
	"github.com/redis/go-redis/v9"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	// 换内存，还是换别的
	// memory 直接基于内存实现
	return memory.NewService()
}
