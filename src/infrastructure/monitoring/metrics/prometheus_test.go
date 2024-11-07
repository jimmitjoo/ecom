package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestProductOperationsMetrics(t *testing.T) {
	// Reset metrics by creating new ones
	ProductOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_product_operations_total",
			Help: "Total number of product operations",
		},
		[]string{"operation", "status"},
	)

	testCases := []struct {
		operation string
		status    string
		value     int
	}{
		{"create", "success", 1},
		{"update", "success", 2},
		{"delete", "failure", 1},
	}

	for _, tc := range testCases {
		t.Run(tc.operation+"_"+tc.status, func(t *testing.T) {
			for i := 0; i < tc.value; i++ {
				ProductOperations.WithLabelValues(tc.operation, tc.status).Inc()
			}

			counter := ProductOperations.WithLabelValues(tc.operation, tc.status)
			value := testutil.ToFloat64(counter)
			assert.Equal(t, float64(tc.value), value)
		})
	}
}

func TestBatchOperationSizeMetrics(t *testing.T) {
	// Create a new histogram for the test
	BatchOperationSize = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "test_batch_operation_size",
			Help:    "Size of batch operations",
			Buckets: []float64{1, 5, 10, 50, 100, 500, 1000},
		},
	)

	testSizes := []float64{1, 5, 50, 100, 500}
	expectedCount := uint64(len(testSizes))
	expectedSum := float64(0)

	for _, size := range testSizes {
		BatchOperationSize.Observe(size)
		expectedSum += size
	}

	// Get metric values
	metric := &dto.Metric{}
	BatchOperationSize.(prometheus.Histogram).Write(metric)

	// Verify count and sum
	assert.Equal(t, expectedCount, metric.Histogram.GetSampleCount())
	assert.Equal(t, expectedSum, metric.Histogram.GetSampleSum())
}

func TestWebSocketConnectionMetrics(t *testing.T) {
	// Create a new gauge for the test
	ActiveWebSocketConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "test_active_websocket_connections",
			Help: "Number of active WebSocket connections",
		},
	)

	testCases := []struct {
		delta    float64
		expected float64
	}{
		{1, 1},
		{5, 6},
		{-2, 4},
	}

	for _, tc := range testCases {
		if tc.delta > 0 {
			ActiveWebSocketConnections.Add(tc.delta)
		} else {
			ActiveWebSocketConnections.Sub(-tc.delta)
		}

		value := testutil.ToFloat64(ActiveWebSocketConnections)
		assert.Equal(t, tc.expected, value)
	}
}

func TestEventProcessingDurationMetrics(t *testing.T) {
	// Create a new histogram for the test
	EventProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "test_event_processing_duration_seconds",
			Help:    "Time spent processing events",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"event_type"},
	)

	testCases := []struct {
		eventType string
		duration  float64
	}{
		{"product.created", 0.1},
		{"product.updated", 0.2},
		{"product.deleted", 0.15},
	}

	for _, tc := range testCases {
		t.Run(tc.eventType, func(t *testing.T) {
			EventProcessingDuration.WithLabelValues(tc.eventType).Observe(tc.duration)

			// Get metric values
			metric := &dto.Metric{}
			observer, err := EventProcessingDuration.GetMetricWithLabelValues(tc.eventType)
			assert.NoError(t, err)
			observer.(prometheus.Histogram).Write(metric)

			// Verify that the value was registered
			assert.Equal(t, uint64(1), metric.Histogram.GetSampleCount())
			assert.Equal(t, tc.duration, metric.Histogram.GetSampleSum())
		})
	}
}

func TestMetricsRegistration(t *testing.T) {
	registry := prometheus.NewRegistry()

	metrics := []prometheus.Collector{
		ProductOperations,
		BatchOperationSize,
		ActiveWebSocketConnections,
		EventProcessingDuration,
		RepositoryOperationDuration,
	}

	for _, metric := range metrics {
		err := registry.Register(metric)
		if err != nil {
			// Ignore "already exists" error since metrics are created by promauto
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				t.Errorf("Failed to register metric: %v", err)
			}
		}
	}
}
