package events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

const topicSyncData = "sync_data_event"

type SaramaSyncProducer struct {
	client sarama.SyncProducer
}

func NewSaramaProducer(client sarama.Client) (*SaramaSyncProducer, error) {
	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}
	return &SaramaSyncProducer{
		p,
	}, nil
}

func (p *SaramaSyncProducer) ProduceSyncEvent(ctx context.Context, evt BizTags) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = p.client.SendMessage(&sarama.ProducerMessage{
		Topic: topicSyncData,
		Value: sarama.ByteEncoder(data),
	})
	return err
}

func (p *SaramaSyncProducer) ProduceStandardSyncEvent(ctx context.Context, tags BizTags) error {
	tdata, _ := json.Marshal(tags)
	evt := SyncDataEvent{
		IndexName: "tags_index",
		// 构成一个唯一的 doc id
		// 要确保后面打了新标签的时候，搜索那边也会有对应的修改
		DocID: fmt.Sprintf("%d_%s_%d", tags.Uid, tags.Biz, tags.BizId),
		Data:  string(tdata),
	}
	data, _ := json.Marshal(evt)
	_, _, err := p.client.SendMessage(&sarama.ProducerMessage{
		Topic: topicSyncData,
		Value: sarama.ByteEncoder(data),
	})
	return err
}
