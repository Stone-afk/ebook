package ginx

import "github.com/prometheus/client_golang/prometheus"

// 包变量导致我们这个地方的代码非常垃圾
var vector *prometheus.CounterVec

func InitCounter(opt prometheus.CounterOpts) {
	vector = prometheus.NewCounterVec(opt, []string{"code"})
	prometheus.MustRegister(vector)
}
