package main

import (
	"go-service-demo/pkg/router"
	"net/http"
)

func main() {
	router := router.InitRouter()
	http.ListenAndServe(":8080", router)
}
