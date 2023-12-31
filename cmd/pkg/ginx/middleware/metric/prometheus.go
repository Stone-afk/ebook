package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

// MiddlewareBuilder 主要是统计响应时间
type MiddlewareBuilder struct {
	// 除了一个 Name 是必选的，其它都是可选的
	// 如果暴露 New 方法，那么就需要考虑暴露其他的方法来允许用户配置 Namespace 等
	// 所以我直接做成了公开字段
	Namespace string
	Subsystem string
	Name      string
	Help      string
	// 这一个实例名字，可以考虑使用 本地 IP，
	// 又或者在启动的时候配置一个 ID
	InstanceID string
}

func (m *MiddlewareBuilder) BuildResponseTime() gin.HandlerFunc {
	// pattern 是指你命中的路由
	// 是指你的 HTTP 的 status
	// path /detail/1
	labels := []string{"method", "pattern", "status"}
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      m.Name + "_resp_time",
		Help:      m.Help,
		ConstLabels: map[string]string{
			"instance_id": m.InstanceID,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.005,
			0.999: 0.0001,
		},
	}, labels)
	prometheus.MustRegister(summary)
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      m.Name + "_active_req",
		Help:      m.Help,
		ConstLabels: map[string]string{
			"instance_id": m.InstanceID,
		},
	})
	prometheus.MustRegister(gauge)
	return func(ctx *gin.Context) {
		start := time.Now()
		gauge.Inc()
		defer func() {
			duration := time.Since(start)
			gauge.Dec()
			// 404????
			pattern := ctx.FullPath()
			if pattern == "" {
				pattern = "unknown"
			}
			summary.WithLabelValues(
				ctx.Request.Method,
				pattern,
				strconv.Itoa(ctx.Writer.Status()),
			).Observe(float64(duration.Milliseconds()))
		}()
		// 你最终就会执行到业务里面
		ctx.Next()
	}
}

func (m *MiddlewareBuilder) BuildActiveRequest() gin.HandlerFunc {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: m.Namespace,
		Subsystem: m.Subsystem,
		Name:      m.Name + "_active_req",
		Help:      m.Help,
		ConstLabels: map[string]string{
			"instance_id": m.InstanceID,
		},
	})
	prometheus.MustRegister(gauge)
	return func(ctx *gin.Context) {
		gauge.Inc()
		defer gauge.Dec()
		ctx.Next()
	}
}
