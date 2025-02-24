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
const UpdateUserUrl = "/{id}"
const DeleteUserUrl = "/{id}"

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

	// Subrouter cho /api/v1
	baseRouter := router.PathPrefix(ApiV1Prefix).Subrouter()
	baseRouter.Use()

	// Subrouter cho /users
	userRouter := baseRouter.PathPrefix(UserControllerPrefix).Subrouter()
	userRouter.Use(
		middleware.NewTrackingMiddleware().Do,
		middleware.NewMonitorMiddleware().Do,
	)
	userController := user_controller.NewUserController(sqlDb, redis)
	userRouter.HandleFunc(FindUserByIdUrl, userController.FindUserById).Methods(constant.GetMethod)
	userRouter.HandleFunc(FindAllUsersUrl, userController.FindAllUsers).Methods(constant.GetMethod)
	userRouter.HandleFunc(CreateNewUserUrl, userController.SaveUser).Methods(constant.PostMethod)
	userRouter.HandleFunc(DeleteUserUrl, userController.DeleteUser).Methods(constant.DeleteMethod)
	userRouter.HandleFunc(UpdateUserUrl, userController.UpdateUser).Methods(constant.PutMethod)
	return router
}
