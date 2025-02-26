package router

import (
	"go-service-demo/internal/drivers/auth_controller"
	"go-service-demo/internal/drivers/user_controller"
	"go-service-demo/internal/middleware"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/database/mysql"
	"go-service-demo/pkg/database/redis"
	"go-service-demo/pkg/utils"
	"log"

	"github.com/gorilla/mux"
)

const ApiV1Prefix = "/api/v1"

const UserControllerPrefix = "/users"
const FindUserByIdUrl = "/{id}"
const FindAllUsersUrl = ""
const UpdateUserUrl = "/{id}"
const DeleteUserUrl = "/{id}"

const AuthControllerPrefix = "/auth"
const LoginUrl = "/login"
const RegisterUrl = "/register"
const LogoutUrl = "/logout"

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

	jwt := utils.NewJwt(
		"token_secret",
		"refresh_token_secret",
		300,
		7*86400,
	)

	router := mux.NewRouter()

	// Subrouter cho /auth
	authRouter := router.PathPrefix(AuthControllerPrefix).Subrouter()
	authRouter.Use()
	authController := auth_controller.NewAuthController(sqlDb, jwt)
	authRouter.HandleFunc(LoginUrl, authController.Login).Methods(constant.PostMethod)
	authRouter.HandleFunc(RegisterUrl, authController.Register).Methods(constant.PostMethod)
	authRouter.HandleFunc(LogoutUrl, authController.Logout).Methods(constant.PostMethod)

	// Subrouter cho /api/v1
	baseRouter := router.PathPrefix(ApiV1Prefix).Subrouter()
	baseRouter.Use(
		middleware.NewJwtMiddleware().Do,
	)

	// Subrouter cho /api/v1/users
	userRouter := baseRouter.PathPrefix(UserControllerPrefix).Subrouter()
	userRouter.Use(
		middleware.NewTrackingMiddleware().Do,
		middleware.NewMonitorMiddleware().Do,
	)
	userController := user_controller.NewUserController(sqlDb, redis)
	userRouter.HandleFunc(FindUserByIdUrl, userController.FindUserById).Methods(constant.GetMethod)
	userRouter.HandleFunc(FindAllUsersUrl, userController.FindAllUsers).Methods(constant.GetMethod)
	userRouter.HandleFunc(DeleteUserUrl, userController.DeleteUser).Methods(constant.DeleteMethod)
	userRouter.HandleFunc(UpdateUserUrl, userController.UpdateUser).Methods(constant.PutMethod)
	return router
}
