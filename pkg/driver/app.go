package driver

import (
	"fmt"
	"go-service-demo/internal/drivers/auth_controller"
	"go-service-demo/internal/drivers/user_controller"
	"go-service-demo/internal/middleware"
	"go-service-demo/pkg"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/database/mysql"
	"go-service-demo/pkg/database/redis"
	"go-service-demo/pkg/messaging_system"
	"go-service-demo/pkg/utils"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

const apiV1Prefix = "/api/v1"

const userControllerPrefix = "/users"
const getUser = ""
const getUsers = "/all"
const updateUser = ""
const deleteUser = ""

const authControllerPrefix = "/auth"
const loginUrl = "/login"
const registerUrl = "/register"
const logoutUrl = "/logout"
const refreshTokenUrl = "/refresh_token"
const verifyUserUrl = "/verify_user"

func Run() {
	pkg.LoadConfig()
	sqlDb := mysql.NewMySql()
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

	jwtAccessTokenTtl, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_TTL_S"))
	jwtRefreshTokenTtl, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_TTL_S"))
	jwt := utils.NewJwt(
		os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
		os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
		jwtAccessTokenTtl,
		jwtRefreshTokenTtl,
	)

	rabbitMq := &messaging_system.RabbitMQ{
		Url:      os.Getenv("RABBITMQ_URL"),
		Protocol: os.Getenv("RABBITMQ_PROTOCOL"),
		Username: os.Getenv("RABBITMQ_USERNAME"),
		Password: os.Getenv("RABBITMQ_PASSWORD"),
	}
	err = rabbitMq.Connect()
	if err != nil {
		log.Println("Error when connect to rabbitmq: " + err.Error())
		panic(err)
	}
	log.Println("Connect to rabbitmq successfully")

	router := mux.NewRouter()

	// Middleware cho toàn bộ router
	router.Use(
		middleware.XssProtectionMiddleware,
		middleware.CorsMiddleware,
		middleware.MonitorMiddleware,
		middleware.ErrorHandlerMiddleware,
	)

	// Subrouter cho /auth
	authRouter := router.PathPrefix(authControllerPrefix).Subrouter()
	authRouter.Use()
	authController := auth_controller.NewAuthController(sqlDb, jwt, rabbitMq)
	authRouter.HandleFunc(loginUrl, authController.Login).Methods(constant.PostMethod)
	authRouter.HandleFunc(registerUrl, authController.Register).Methods(constant.PostMethod)
	authRouter.HandleFunc(logoutUrl, authController.Logout).Methods(constant.PostMethod)
	authRouter.HandleFunc(refreshTokenUrl, authController.RefreshToken).Methods(constant.PostMethod)
	authRouter.HandleFunc(verifyUserUrl, authController.VerifyUser).Methods(constant.GetMethod)

	// Subrouter cho /api/v1
	baseRouter := router.PathPrefix(apiV1Prefix).Subrouter()
	baseRouter.Use(
		middleware.NewJwtMiddleware(jwt).Do,
	)

	// Subrouter cho /api/v1/users
	userRouter := baseRouter.PathPrefix(userControllerPrefix).Subrouter()
	userController := user_controller.NewUserController(sqlDb, redis)
	userRouter.HandleFunc(getUser, userController.GetUser).Methods(constant.GetMethod)
	userRouter.HandleFunc(getUsers, userController.GetUsers).Methods(constant.GetMethod)
	userRouter.HandleFunc(deleteUser, userController.DeleteUser).Methods(constant.DeleteMethod)
	userRouter.HandleFunc(updateUser, userController.UpdateUser).Methods(constant.PutMethod)

	fmt.Println("Server is running on port: " + os.Getenv("SERVER_PORT"))
	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), router)
}
