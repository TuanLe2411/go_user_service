package middleware

import (
	"go-service-demo/pkg/constant"
	"net/http"
)

type JwtMiddleware struct {
}

func NewJwtMiddleware() constant.Middleware {
	return &JwtMiddleware{}
}

func (j *JwtMiddleware) Do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
