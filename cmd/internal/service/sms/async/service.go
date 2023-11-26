package async

import "ebook/cmd/internal/service/sms"

type Service struct {
	svc sms.Service
	// 转异步，存储发短信请求的 repository
}
