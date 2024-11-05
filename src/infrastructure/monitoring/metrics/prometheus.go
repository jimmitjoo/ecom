package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Product operations metrics
	ProductOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "product_operations_total",
			Help: "Total number of product operations",
		},
		[]string{"operation", "status"},
	)

	// Batch operations metrics
	BatchOperationSize = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "batch_operation_size",
			Help:    "Size of batch operations",
			Buckets: []float64{1, 5, 10, 50, 100, 500, 1000},
		},
	)

	// WebSocket metrics
	ActiveWebSocketConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_websocket_connections",
			Help: "Number of active WebSocket connections",
		},
	)

	// Event processing metrics
	EventProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "event_processing_duration_seconds",
			Help:    "Time spent processing events",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"event_type"},
	)

	// Repository operation latency
	RepositoryOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "repository_operation_duration_seconds",
			Help:    "Time spent on repository operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)
