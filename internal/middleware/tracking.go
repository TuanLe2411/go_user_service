package middleware

import (
	"fmt"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/database"
	"math/rand/v2"
	"net/http"
)

type TrackingMiddleware struct {
	db database.Database
}

func NewTrackingMiddleware(db database.Database) constant.Middleware {
	return &TrackingMiddleware{
		db: db,
	}
}

func (t *TrackingMiddleware) Do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trackingId := fmt.Sprintf("%d", rand.Int())
		r.Header.Set("X-Tracking-Id", trackingId)
		w.Header().Set("X-Tracking-Id", trackingId)
		fmt.Printf("TrackingId: %s, request: %s %s \n", trackingId, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
