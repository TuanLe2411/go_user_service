package main

import (
	"go-service-demo/pkg/database/mysql"
	"go-service-demo/pkg/router"
	"log"
	"net/http"
)

func main() {
	sqlDb := mysql.NewMySql("root:root@tcp(localhost:3306)/go_service_demo?parseTime=true&loc=Local&charset=utf8mb4")
	err := sqlDb.Connect()
	if err != nil {
		log.Println("Error when connect to db: " + err.Error())
		panic(err)
	}
	err = sqlDb.Ping()
	if err != nil {
		log.Println("Error when ping to db: " + err.Error())
		panic(err)
	}
	log.Println("Connect to db successfully")
	router := router.InitRouter(sqlDb)
	http.ListenAndServe(":8080", router)
}
