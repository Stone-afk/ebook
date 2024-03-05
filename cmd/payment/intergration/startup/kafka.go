package startup

import (
	"ebook/cmd/payment/events"
	"github.com/IBM/sarama"
)

func InitKafka() sarama.Client {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
	client, err := sarama.NewClient([]string{"localhost:19094"}, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

func NewSyncProducer(client sarama.Client) events.Producer {
	res, err := events.NewSaramaProducer(client)
	if err != nil {
		panic(err)
	}
	return res
}
