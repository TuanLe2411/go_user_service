package main

import (
	"fmt"
	"go-service-demo/internal/handler"
	"go-service-demo/internal/middleware"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/database"
	"go-service-demo/pkg/database/mock"
	"go-service-demo/pkg/util"
	"net/http"
)

func Init(mux *http.ServeMux, db database.Database) {
	mux.Handle(constant.GET_METHOD+constant.API_V1_PREFIX+constant.TEST_URL,
		util.ChainMiddlewares(
			handler.NewGetTestHandler(db),
			middleware.NewTrackingMiddleware(db),
		),
	)
}

func main() {
	mux := http.NewServeMux()
	databaseMock := mock.NewDatabaseMock()
	err := databaseMock.Connect()
	if err != nil {
		fmt.Println("Error when connect to db: " + err.Error())
		panic(err)
	}
	Init(mux, databaseMock)
	http.ListenAndServe(":8080", mux)
}
