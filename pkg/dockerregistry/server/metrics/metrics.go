package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	registryNamespace = "openshift"
	registrySubsystem = "registry"
)

var (
	RegistryAPIRequests *prometheus.SummaryVec
)

// Register the metrics.
func Register() {
	RegistryAPIRequests = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: registryNamespace,
			Subsystem: registrySubsystem,
			Name:      "request_duration_seconds",
			Help:      "Request latency summary in microseconds for each operation",
		},
		[]string{"operation", "name"},
	)
	prometheus.MustRegister(RegistryAPIRequests)
}

// Timer is a helper type to time functions.
type Timer interface {
	// Stop records the duration passed since the Timer was created with NewTimer.
	Stop()
}

// NewTimer wraps the SummaryVec and used to track amount of time passed since the Timer was created.
func NewTimer(collector *prometheus.SummaryVec, labels []string) Timer {
	return &metricTimer{
		collector: collector,
		labels:    labels,
		startTime: time.Now(),
	}
}

type metricTimer struct {
	collector *prometheus.SummaryVec
	labels    []string
	startTime time.Time
}

func (m *metricTimer) Stop() {
	m.collector.WithLabelValues(m.labels...).Observe(float64(time.Since(m.startTime) / time.Second))
}
