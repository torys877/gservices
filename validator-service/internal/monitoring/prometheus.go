package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"endpoint"},
	)
	ResponseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_time_seconds",
			Help:    "Response time distribution",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
)

func InitPrometheus() {
	prometheus.MustRegister(TotalRequests, ResponseDuration)
}
