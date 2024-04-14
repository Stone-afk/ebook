package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const topicFollowerFeedEvent = "follower_feed_event"
const topicFeedEvent = "feed_event"

type FollowerFeedEventProducer struct {
	producer sarama.SyncProducer
}

func (p *FollowerFeedEventProducer) ProduceStandardFeedEvent(ctx context.Context, evt FeedEvent) error {
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

func (p *FollowerFeedEventProducer) ProduceFeedEvent(ctx context.Context, evt FeedEvent) error {
	follower, err := evt.Metadata.Get("follower").AsInt64()
	if err != nil {
		return err
	}
	followee, err := evt.Metadata.Get("followee").AsInt64()
	if err != nil {
		return err
	}
	bizId, err := evt.Metadata.Get("biz_id").AsInt64()
	if err != nil {
		return err
	}
	//biz, err := evt.Metadata.Get("biz").String()
	//if err != nil {
	//	return err
	//}
	val, err := json.Marshal(FollowerFeedEvent{
		Follower: follower,
		Followee: followee,
		Biz:      "Follower",
		BizId:    bizId,
	})
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicFollowerFeedEvent,
		Value: sarama.ByteEncoder(val),
	})
	return err
}

func NewFollowerFeedEventProducer(pc sarama.SyncProducer) FeedEventProducer {
	return &FollowerFeedEventProducer{
		producer: pc,
	}
}
