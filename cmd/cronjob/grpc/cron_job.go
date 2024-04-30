package grpc

import (
	"context"
	cronjobv1 "ebook/cmd/api/proto/gen/cronjob/v1"
	"ebook/cmd/cronjob/domain"
	"ebook/cmd/cronjob/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CronJobServiceServer struct {
	svc service.CronJobService
	cronjobv1.UnimplementedCronJobServiceServer
}

func NewCronJobServiceServer(svc service.CronJobService) *CronJobServiceServer {
	return &CronJobServiceServer{
		svc: svc,
	}
}

func (s *CronJobServiceServer) Register(server grpc.ServiceRegistrar) {
	cronjobv1.RegisterCronJobServiceServer(server, s)
}

func (s *CronJobServiceServer) Preempt(ctx context.Context, request *cronjobv1.PreemptRequest) (*cronjobv1.PreemptResponse, error) {
	job, err := s.svc.Preempt(ctx)
	return &cronjobv1.PreemptResponse{
		Job: convertToV(job),
	}, err
}

func (s *CronJobServiceServer) ResetNextTime(ctx context.Context, request *cronjobv1.ResetNextTimeRequest) (*cronjobv1.ResetNextTimeResponse, error) {
	err := s.svc.ResetNextTime(ctx, convertToDomain(request.Job))
	return &cronjobv1.ResetNextTimeResponse{}, err
}

func (s *CronJobServiceServer) AddJob(ctx context.Context, request *cronjobv1.AddJobRequest) (*cronjobv1.AddJobResponse, error) {
	err := s.svc.AddJob(ctx, convertToDomain(request.Job))
	return &cronjobv1.AddJobResponse{}, err
}

func convertToV(domainCronJob domain.CronJob) *cronjobv1.CronJob {
	return &cronjobv1.CronJob{
		Id:         domainCronJob.Id,
		Name:       domainCronJob.Name,
		Executor:   domainCronJob.Executor,
		Cfg:        domainCronJob.Cfg,
		Expression: domainCronJob.Expression,
		NextTime:   timestamppb.New(domainCronJob.NextTime),
	}
}

func convertToDomain(vCronJob *cronjobv1.CronJob) domain.CronJob {
	return domain.CronJob{
		Id:         vCronJob.GetId(),
		Name:       vCronJob.GetName(),
		Executor:   vCronJob.GetExecutor(),
		Cfg:        vCronJob.GetCfg(),
		Expression: vCronJob.GetExpression(),
		NextTime:   vCronJob.GetNextTime().AsTime(),
	}
}
