package events

import (
	"ebook/cmd/im/service"
	"ebook/cmd/pkg/canalx"
	"ebook/cmd/pkg/logger"
	"github.com/IBM/sarama"
)

type MySQLBinlogConsumer struct {
	client sarama.Client
	l      logger.Logger
	svc    service.UserService
}

func NewMySQLBinlogConsumer(client sarama.Client,
	l logger.Logger, svc service.UserService) *MySQLBinlogConsumer {
	return &MySQLBinlogConsumer{
		client: client,
		l:      l,
		svc:    svc,
	}
}

func (r *MySQLBinlogConsumer) Start() error {
	panic("")
}

func (r *MySQLBinlogConsumer) Consume(msg *sarama.ConsumerMessage,
	val canalx.Message[User]) error {
	panic("")
}

type User struct {
	Id            int64
	Email         string
	Password      string
	Phone         string
	Birthday      string
	Nickname      string
	AboutMe       string
	WechatOpenId  string
	WechatUnionId string

	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
