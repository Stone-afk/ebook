package async_sms

import "github.com/ecodeclub/ekit/sqlx"

type SmsConfig struct {
	TplId   string
	Args    []string
	Numbers []string
}

type AsyncSms struct {
	Id int64
	// 使用我在 ekit 里面支持的 JSON 字段
	Config sqlx.JsonColumn[SmsConfig]
	// 重试次数
	RetryCnt int
	// 重试的最大次数
	RetryMax int
	Status   uint8
	Ctime    int64
	Utime    int64 `gorm:"index"`
}
