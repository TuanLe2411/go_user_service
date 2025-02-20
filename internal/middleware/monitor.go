package middleware

import (
	"go-service-demo/pkg/constant"
	"log"
	"net/http"
	"time"
)

type MonitorMiddleware struct{}

func NewMonitorMiddleware() constant.Middleware {
	return &MonitorMiddleware{}
}

func (m *MonitorMiddleware) Do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("Request: %s %s, duration: %s\n", r.Method, r.URL.Path, duration)
	})
}
