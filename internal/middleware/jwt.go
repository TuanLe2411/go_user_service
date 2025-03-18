package middleware

import (
	"errors"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/utils"
	"net/http"
	"strconv"
)

type JwtMiddleware struct {
	*utils.Jwt
}

func NewJwtMiddleware(jwt *utils.Jwt) constant.Middleware {
	return &JwtMiddleware{
		jwt,
	}
}

func (j *JwtMiddleware) Do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			utils.SetHttpReponseError(r, utils.ErrUnAuthorized, errors.New("missing jwt token"))
			return
		}
		token = token[len("Bearer "):]
		isValid, claims := j.ValidateToken(token)
		if !isValid {
			utils.SetHttpReponseError(r, utils.ErrUnAuthorized, errors.New("invalid jwt token"))
			return
		}
		r.Header.Set("user_id", strconv.FormatInt(claims.UserId, 10))
		r.Header.Set("username", claims.Username)
		next.ServeHTTP(w, r)
	})
}
