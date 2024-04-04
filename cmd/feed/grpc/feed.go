package grpc

import (
	"context"
	feedv1 "ebook/cmd/api/proto/gen/feed/v1"
	"ebook/cmd/feed/domain"
	"ebook/cmd/feed/service"
	"encoding/json"
	"google.golang.org/grpc"
	"time"
)

type FeedEventServiceServer struct {
	feedv1.UnimplementedFeedSvcServer
	svc service.FeedService
}

func NewFeedEventServiceServer(svc service.FeedService) *FeedEventServiceServer {
	return &FeedEventServiceServer{
		svc: svc,
	}
}

func (s *FeedEventServiceServer) CreateFeedEvent(ctx context.Context, request *feedv1.CreateFeedEventRequest) (*feedv1.CreateFeedEventResponse, error) {
	err := s.svc.CreateFeedEvent(ctx, s.convertToDomain(request.GetFeedEvent()))
	return &feedv1.CreateFeedEventResponse{}, err
}

func (s *FeedEventServiceServer) FindFeedEvents(ctx context.Context, request *feedv1.FindFeedEventsRequest) (*feedv1.FindFeedEventsResponse, error) {
	eventList, err := s.svc.GetFeedEvents(ctx, request.GetUid(), request.Timestamp, request.Limit)
	if err != nil {
		return &feedv1.FindFeedEventsResponse{}, err
	}
	res := make([]*feedv1.FeedEvent, 0, len(eventList))
	for _, event := range eventList {
		res = append(res, s.convertToView(event))
	}
	return &feedv1.FindFeedEventsResponse{
		FeedEvents: res,
	}, nil
}

func (s *FeedEventServiceServer) Register(server grpc.ServiceRegistrar) {
	feedv1.RegisterFeedSvcServer(server, s)
}

func (s *FeedEventServiceServer) convertToDomain(event *feedv1.FeedEvent) domain.FeedEvent {
	ext := map[string]string{}
	_ = json.Unmarshal([]byte(event.Content), &ext)
	return domain.FeedEvent{
		ID:    event.Id,
		Ctime: time.Unix(event.Ctime, 0),
		Type:  event.GetType(),
		Ext:   ext,
	}
}

func (s *FeedEventServiceServer) convertToView(event domain.FeedEvent) *feedv1.FeedEvent {
	val, _ := json.Marshal(event.Ext)
	return &feedv1.FeedEvent{
		Id:      event.ID,
		Type:    event.Type,
		Ctime:   event.Ctime.Unix(),
		Content: string(val),
	}
}
