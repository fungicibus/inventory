package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

func registerMetrics(appVersion string, router *chi.Mux) {
	var (
		version = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "version",
			Help: "Version information about this binary",
			ConstLabels: map[string]string{
				"version": appVersion,
			},
		})
		httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all HTTP requests",
		}, []string{"code", "method"})

		httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of all HTTP requests",
		}, []string{"code", "handler", "method"})
	)
	version.Set(1)

	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(version)

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if skipMonitoring(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			statusCode := fmt.Sprintf("%d", ww.Status())

			httpRequestsTotal.WithLabelValues(statusCode, r.Method).Inc()
			httpRequestDuration.WithLabelValues(
				statusCode,
				r.URL.Path,
				r.Method,
			).Observe(duration.Seconds())
		})
	})
}

func skipMonitoring(urlPath string) bool {
	return (urlPath == "/metrics" ||
		urlPath == "/healthcheck" ||
		strings.Contains(urlPath, "swagger"))
}
