package app_controller

import "net/http"

type AppController struct{}

func (a *AppController) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
