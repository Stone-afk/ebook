package events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type BizTags struct {
	Uid   int64  `json:"uid"`
	Biz   string `json:"biz"`
	BizId int64  `json:"biz_id"`
	// 只传递 string
	Tags []string `json:"tags"`
}

type SyncDataEvent struct {
	IndexName string
	DocID     string
	Data      string
}

type Producer interface {
	ProduceSyncEvent(ctx context.Context, data BizTags) error
}

type SaramaSyncProducer struct {
	client sarama.SyncProducer
}

func (p *SaramaSyncProducer) ProduceSyncEvent(ctx context.Context, tags BizTags) error {
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
		Topic: "search_sync_data",
		Value: sarama.ByteEncoder(data),
	})
	return err
}
