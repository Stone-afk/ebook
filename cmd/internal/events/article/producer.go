package article

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const topicReadEvent = "article_read_event"

type Producer interface {
	ProduceReadEvent(ctx context.Context, evt ReadEvent) error
	//ProduceReadEventV1(ctx context.Context, v1 ReadEventV1) error
}

type ReadEvent struct {
	Uid int64
	Aid int64
}

type ReadEventV1 struct {
	Uids []int64
	Aids []int64
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(pc sarama.SyncProducer) Producer {
	return &KafkaProducer{
		producer: pc,
	}
}

//func (k *KafkaProducer) ProduceReadEventV1(ctx context.Context, v1 ReadEventV1) error {
//	//TODO implement me
//	panic("implement me")
//}

func (k *KafkaProducer) ProduceReadEvent(ctx context.Context, evt ReadEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicReadEvent,
		Value: sarama.ByteEncoder(val),
	})
	return err
}
