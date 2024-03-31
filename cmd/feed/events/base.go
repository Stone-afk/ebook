package events

const topicFeedEvent = "feed_event"

type FeedEvent struct {
	Type     string
	Metadata map[string]string
}
