package middleware

import (
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/utils"
	"net/http"
)

type JwtMiddleware struct {
	jwt *utils.Jwt
}

func NewJwtMiddleware(jwt *utils.Jwt) constant.Middleware {
	return &JwtMiddleware{
		jwt: jwt,
	}
}

func (j *JwtMiddleware) Do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from header
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token = token[len("Bearer "):]
		isValid, claims := j.jwt.ValidateToken(token)
		if !isValid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		r.Header.Set("username", claims.Username)
		next.ServeHTTP(w, r)
	})
}
