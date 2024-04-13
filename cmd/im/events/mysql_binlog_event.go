package events

import (
	"ebook/cmd/im/service"
	"ebook/cmd/pkg/logger"
	"github.com/IBM/sarama"
)

type MySQLBinlogConsumer struct {
	client sarama.Client
	l      logger.Logger
	svc    service.UserService
}
