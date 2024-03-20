package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"strconv"
)

const topicSyncData = "sync_data_event"

const topicSyncUser = "sync_user_event"

type SaramaSyncProducer struct {
	client sarama.SyncProducer
}

func (p *SaramaSyncProducer) ProduceSyncEvent(ctx context.Context, evt UserEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = p.client.SendMessage(&sarama.ProducerMessage{
		Topic: topicSyncUser,
		Value: sarama.ByteEncoder(data),
	})
	return err
}

func (p *SaramaSyncProducer) ProduceStandardSyncEvent(ctx context.Context, ue UserEvent) error {
	tdata, err := json.Marshal(ue)
	if err != nil {
		return err
	}
	evt := SyncDataEvent{
		IndexName: "user_index",
		// 构成一个唯一的 doc id
		// 要确保后面打了新标签的时候，搜索那边也会有对应的修改
		DocID: strconv.FormatInt(ue.Id, 10),
		Data:  string(tdata),
	}
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

func NewSaramaSyncProducer(pc sarama.SyncProducer) SyncSearchEventProducer {
	return &SaramaSyncProducer{
		pc,
	}
}
