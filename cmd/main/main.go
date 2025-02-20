package main

import (
	"go-service-demo/internal/handler"
	"go-service-demo/internal/middleware"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/database"
	"go-service-demo/pkg/database/sql_db_mock"
	"go-service-demo/pkg/util"
	"log"
	"net/http"
)

func Init(mux *http.ServeMux, db database.SqlDatabase) {
	mux.Handle(constant.GET_METHOD+constant.API_V1_PREFIX+constant.TEST_URL,
		util.ChainMiddlewares(
			handler.NewGetTestHandler(db),
			middleware.NewTrackingMiddleware(),
			middleware.NewMonitorMiddleware(),
		),
	)
}

func main() {
	mux := http.NewServeMux()
	databaseMock := sql_db_mock.NewSqlDatabaseMock()
	err := databaseMock.Connect()
	if err != nil {
		log.Println("Error when connect to db: " + err.Error())
		panic(err)
	}
	Init(mux, databaseMock)
	http.ListenAndServe(":8080", mux)
}
