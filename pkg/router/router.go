package router

import (
	"go-service-demo/internal/controllers/user_controller"
	"go-service-demo/internal/middleware"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/database/mysql"
	"go-service-demo/pkg/database/redis"
	"log"

	"github.com/gorilla/mux"
)

const ApiV1Prefix = "/api/v1"

const UserControllerPrefix = "/users"

const FindUserByIdUrl = "/{id}"
const FindAllUsersUrl = ""
const CreateNewUserUrl = ""

func InitRouter() *mux.Router {
	sqlDb := mysql.NewMySql("root:root@tcp(localhost:3306)/go_service_demo?parseTime=true&loc=Local&charset=utf8mb4")
	err := sqlDb.Connect()
	if err != nil {
		log.Println("Error when connect to db: " + err.Error())
		panic(err)
	}
	log.Println("Connect to db successfully")

	redis := redis.NewRedisClient()
	err = redis.Connect()
	if err != nil {
		log.Println("Error when connect to redis: " + err.Error())
		panic(err)
	}
	log.Println("Connect to redis successfully")

	router := mux.NewRouter()

	// Subrouter cho /users
	userRouter := router.PathPrefix(ApiV1Prefix + UserControllerPrefix).Subrouter()
	userRouter.Use(
		middleware.NewTrackingMiddleware().Do,
		middleware.NewMonitorMiddleware().Do,
	)
	userRouter.HandleFunc(FindUserByIdUrl, user_controller.NewUserController(sqlDb, redis).FindUserById).Methods(constant.GetMethod)
	userRouter.HandleFunc(FindAllUsersUrl, user_controller.NewUserController(sqlDb, redis).FindAllUsers).Methods(constant.GetMethod)
	userRouter.HandleFunc(CreateNewUserUrl, user_controller.NewUserController(sqlDb, redis).SaveUser).Methods(constant.PostMethod)

	return router
}
