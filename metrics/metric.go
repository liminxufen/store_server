package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ServerTag = "store_server"
)

var RequestTotalCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Subsystem: "request_total",
	Name:      "counter",
	Help:      "total number of http request received",
}, []string{ServerTag})

var RequestMethodCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Subsystem: "request_method",
	Name:      "counter",
	Help:      "number of http request method received",
}, []string{ServerTag, "method"})

var RequestStatusCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Subsystem: "request_status",
	Name:      "counter",
	Help:      "number of http request status received",
}, []string{ServerTag, "status"})

var RequestClassifyTotalCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Subsystem: "request_classify",
	Name:      "counter",
	Help:      "classify total desc of http request received",
}, []string{ServerTag, "status", "host", "method", "url"})

var RequestHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Subsystem: "request_histogram",
	Name:      "latency",
	Help:      "Latency of request in seconds.",
	Buckets:   prometheus.ExponentialBuckets(0.0005, 2, 22), // ~ 17min
}, []string{ServerTag, "method", "status"})

var RequestSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Subsystem: "request_summary",
	Name:      "latency_summary",
	Help:      "Latency of request in seconds.",
	Objectives: map[float64]float64{0.5: 0.05, 0.8: 0.02, 0.9: 0.01, 0.93: 0.01,
		0.96: 0.001, 0.99: 0.001, 1: 0.0001}, // ~ 8s
}, []string{ServerTag, "method"})

var RequestClassifySummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Subsystem: "request_classify_summary",
	Name:      "latency_summary",
	Help:      "Latency of request in seconds.",
	Objectives: map[float64]float64{0.5: 0.05, 0.8: 0.02, 0.9: 0.01, 0.93: 0.01,
		0.96: 0.001, 0.99: 0.001, 1: 0.0001}, // ~ 8s
}, []string{ServerTag, "status", "host", "method", "url"})

var AlbumPublishCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Subsystem: "album_publish",
	Name:      "counter",
	Help:      "total desc of auto published albums",
}, []string{ServerTag, "type", "subtype"})

func init() {
	prometheus.MustRegister(
		RequestTotalCounter,
		RequestMethodCounter,
		RequestStatusCounter,
		RequestClassifyTotalCounter,
		RequestHistogram,
		RequestSummary,
		RequestClassifySummary,
		AlbumPublishCounter,
	)
}
