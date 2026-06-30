package middlewares

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	http_requests_count = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total count of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)
	http_request_duration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Latency of HTTP requests in seconds",
		},
		[]string{"method", "endpoint", "status_code"},
	)
)

func init() {
	prometheus.MustRegister(http_requests_count)
	prometheus.MustRegister(http_request_duration)
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ww := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(ww, r)

		statusCode := strconv.Itoa(ww.statusCode)

		http_requests_count.WithLabelValues(r.Method, r.URL.Path, statusCode).Inc()
		http_request_duration.WithLabelValues(r.Method, r.URL.Path, statusCode).Observe(time.Since(startTime).Seconds())
	})
}

func MetricHandler() http.Handler {
	return promhttp.Handler()
}
