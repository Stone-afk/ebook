package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

const topicArticleFeedEvent = "article_feed_event"
const topicFeedEvent = "feed_event"

// ArticleFeedEvent 由业务方定义，本服务做适配
type ArticleFeedEvent struct {
	Uid   int64
	Aid   int64
	Title string
}

type ArticleFeedEventProducer struct {
	producer sarama.SyncProducer
}

func (p *ArticleFeedEventProducer) ProduceStandardFeedEvent(ctx context.Context, evt FeedEvent) error {
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

func (p *ArticleFeedEventProducer) ProduceFeedEvent(ctx context.Context, evt FeedEvent) error {
	uid, err := evt.Metadata.Get("uid").AsInt64()
	if err != nil {
		return err
	}
	aid, err := evt.Metadata.Get("aid").AsInt64()
	if err != nil {
		return err
	}
	title, err := evt.Metadata.Get("title").String()
	if err != nil {
		return err
	}
	val, err := json.Marshal(ArticleFeedEvent{
		Uid:   uid,
		Aid:   aid,
		Title: title,
	})
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topicArticleFeedEvent,
		Value: sarama.ByteEncoder(val),
	})
	return err
}

func NewArticleFeedEventProducer(pc sarama.SyncProducer) FeedEventProducer {
	return &ArticleFeedEventProducer{
		producer: pc,
	}
}
