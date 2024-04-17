package metric

import (
	"context"
	"ebook/cmd/sms/service"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type PrometheusService struct {
	svc    service.Service
	vector *prometheus.SummaryVec
}

func NewService(svc service.Service,
	namespace string,
	subsystem string,
	instanceId string,
	name string) *PrometheusService {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      name,
		ConstLabels: map[string]string{
			"instance_id": instanceId,
		},
		Objectives: map[float64]float64{
			0.9:   0.01,
			0.95:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"tpl"})
	prometheus.MustRegister(vector)
	return &PrometheusService{
		vector: vector,
		svc:    svc,
	}
}

func (p *PrometheusService) Send(ctx context.Context, tplId string,
	args []string, numbers ...string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		p.vector.WithLabelValues(tplId).
			Observe(float64(duration.Milliseconds()))
	}()
	return p.svc.Send(ctx, tplId, args, numbers...)
}
