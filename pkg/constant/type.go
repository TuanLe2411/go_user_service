package constant

import "net/http"

type Middleware interface {
	Do(http.Handler) http.Handler
}
