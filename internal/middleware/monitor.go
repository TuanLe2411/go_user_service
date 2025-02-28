package middleware

import (
	"fmt"
	"go-service-demo/pkg/constant"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

type MonitorMiddleware struct{}

func NewMonitorMiddleware() constant.Middleware {
	return &MonitorMiddleware{}
}

func (m *MonitorMiddleware) Do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trackingId := fmt.Sprintf("%d", rand.Int())
		r.Header.Set("X-Tracking-Id", trackingId)
		w.Header().Set("X-Tracking-Id", trackingId)
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("TrackingId: %s, Request: %s %s, duration: %s\n", trackingId, r.Method, r.URL.Path, duration)
	})
}
