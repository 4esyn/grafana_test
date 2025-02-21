package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Метрики времени выполнения запросов к эндпоинтам
	EndpointDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "geoservice_endpoint_duration_seconds",
			Help: "Время выполнения запроса к эндпоинту",
		},
		[]string{"endpoint"},
	)

	// Метрики количества запросов к эндпоинтам
	EndpointRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "geoservice_endpoint_requests_total",
			Help: "Общее количество запросов к эндпоинту",
		},
		[]string{"endpoint"},
	)

	// Метрики времени обращения к кэшу
	CacheAccessDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "geoservice_cache_duration_seconds",
			Help: "Время обращения к кэшу",
		},
		[]string{"method"},
	)

	// Метрики времени обращения к БД
	DBAccessDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "geoservice_db_duration_seconds",
			Help: "Время обращения к БД",
		},
		[]string{"method"},
	)

	// Метрики времени обращения к внешнему API
	ExternalAPIDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "geoservice_external_api_duration_seconds",
			Help: "Время обращения к внешнему API",
		},
		[]string{"method"},
	)
)

// MeasureEndpointDuration измеряет время выполнения эндпоинта
func MeasureEndpointDuration(endpoint string) func() {
	start := time.Now()
	EndpointRequests.WithLabelValues(endpoint).Inc()
	return func() {
		duration := time.Since(start).Seconds()
		EndpointDuration.WithLabelValues(endpoint).Observe(duration)
	}
}

// MeasureCacheAccess измеряет время обращения к кэшу
func MeasureCacheAccess(method string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start).Seconds()
		CacheAccessDuration.WithLabelValues(method).Observe(duration)
	}
}

// MeasureDBAccess измеряет время обращения к БД
func MeasureDBAccess(method string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start).Seconds()
		DBAccessDuration.WithLabelValues(method).Observe(duration)
	}
}

// MeasureExternalAPIAccess измеряет время обращения к внешнему API
func MeasureExternalAPIAccess(method string) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start).Seconds()
		ExternalAPIDuration.WithLabelValues(method).Observe(duration)
	}
}
