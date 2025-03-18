package metrics

import (
	"context"
	"fizzbuzz-server/internal/config"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	// FizzBuzzRequests counts the total number of FizzBuzz requests
	FizzBuzzRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "fizzbuzz_requests_total",
			Help: "The total number of FizzBuzz requests",
		},
		[]string{"status"},
	)

	// FizzBuzzRequestDuration measures the duration of FizzBuzz requests
	FizzBuzzRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "fizzbuzz_request_duration_seconds",
			Help:    "The duration of FizzBuzz requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)

	// FizzBuzzResultSize measures the size of FizzBuzz results
	FizzBuzzResultSize = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "fizzbuzz_result_size",
			Help:    "The size of FizzBuzz results (number of elements)",
			Buckets: []float64{10, 50, 100, 500, 1000, 5000, 10000},
		},
	)

	// StatsRequests counts the total number of Stats requests
	StatsRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "stats_requests_total",
			Help: "The total number of Stats requests",
		},
		[]string{"status"},
	)

	// StatsRequestDuration measures the duration of Stats requests
	StatsRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "stats_request_duration_seconds",
			Help:    "The duration of Stats requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)

	// StatsEntriesCount measures the number of entries in the stats map
	StatsEntriesCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "stats_entries_count",
			Help: "The number of entries in the stats map",
		},
	)

	// MostFrequentRequestHits tracks the hit count of the most frequent request
	MostFrequentRequestHits = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "most_frequent_request_hits",
			Help: "The hit count of the most frequent request",
		},
	)
)

// InitPushGateway initializes the Prometheus push gateway if enabled
func InitPushGateway() {
	cfg := config.Get()
	if !cfg.Prometheus.Enabled || cfg.Prometheus.PushGateway == "" {
		return
	}

	// Start a goroutine to periodically push metrics to the Prometheus push gateway
	go func() {
		pusher := push.New(cfg.Prometheus.PushGateway, "fizzbuzz_server")
		
		// Add all collectors
		registry := prometheus.NewRegistry()
		registry.MustRegister(
			FizzBuzzRequests,
			FizzBuzzRequestDuration,
			FizzBuzzResultSize,
			StatsRequests,
			StatsRequestDuration,
			StatsEntriesCount,
			MostFrequentRequestHits,
		)
		pusher.Gatherer(registry)

		for {
			if err := pusher.Push(); err != nil {
				log.Printf("Error pushing metrics to Prometheus: %v", err)
			}
			time.Sleep(cfg.Prometheus.PushInterval)
		}
	}()
}

// Shutdown gracefully shuts down the metrics system
func Shutdown(ctx context.Context) error {
	// Perform any cleanup needed for metrics
	return nil
}