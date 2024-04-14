package like

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const topicLikeFeedEvent = "liked_feed_event"
const topicFeedEvent = "feed_event"

type LikedFeedEventProducer struct {
	producer sarama.SyncProducer
}

func (p *LikedFeedEventProducer) ProduceStandardFeedEvent(ctx context.Context, evt FeedEvent) error {
	val, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicFeedEvent,
		Value: sarama.ByteEncoder(val),
	})
	return err
}

func (p *LikedFeedEventProducer) ProduceFeedEvent(ctx context.Context, evt FeedEvent) error {
	uid, err := evt.Metadata.Get("uid").AsInt64()
	if err != nil {
		return err
	}
	Liked, err := evt.Metadata.Get("liked").AsInt()
	if err != nil {
		return err
	}
	bizId, err := evt.Metadata.Get("biz_id").AsInt64()
	if err != nil {
		return err
	}
	biz, err := evt.Metadata.Get("biz").String()
	if err != nil {
		return err
	}
	val, err := json.Marshal(LikeFeedEvent{
		Uid:   uid,
		Liked: Liked,
		Biz:   biz,
		BizId: bizId,
	})
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicLikeFeedEvent,
		Value: sarama.ByteEncoder(val),
	})
	return err
}

func NewLikedFeedEventProducer(pc sarama.SyncProducer) FeedEventProducer {
	return &LikedFeedEventProducer{
		producer: pc,
	}
}
