package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const topicSyncOpenIM = "open_im_binlog"

type SyncOpenIMProducer struct {
	client sarama.SyncProducer
}

func (p *SyncOpenIMProducer) ProduceIMEvent(ctx context.Context, evt IMUserEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = p.client.SendMessage(&sarama.ProducerMessage{
		Topic: topicSyncOpenIM,
		Value: sarama.ByteEncoder(data),
	})
	return err
}

func NewSyncOpenIMProducer(pc sarama.SyncProducer) SyncIMEventProducer {
	return &SyncOpenIMProducer{
		pc,
	}
}
