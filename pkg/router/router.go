package router

import (
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
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

const ApiV1Prefix = "/api/v1"

const UserControllerPrefix = "/users"
const GetUser = ""
const GetUsers = "/all"
const UpdateUser = ""
const DeleteUser = ""

const AuthControllerPrefix = "/auth"
const LoginUrl = "/login"
const RegisterUrl = "/register"
const LogoutUrl = "/logout"
const RefreshTokenUrl = "/refresh_token"
const verifyUserUrl = "/verify_user"

func InitRouter() *mux.Router {
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
		middleware.NewMonitorMiddleware().Do,
		middleware.ErrorHandlerMiddleware,
	)

	// Subrouter cho /auth
	authRouter := router.PathPrefix(AuthControllerPrefix).Subrouter()
	authRouter.Use()
	authController := auth_controller.NewAuthController(sqlDb, jwt, rabbitMq)
	authRouter.HandleFunc(LoginUrl, authController.Login).Methods(constant.PostMethod)
	authRouter.HandleFunc(RegisterUrl, authController.Register).Methods(constant.PostMethod)
	authRouter.HandleFunc(LogoutUrl, authController.Logout).Methods(constant.PostMethod)
	authRouter.HandleFunc(RefreshTokenUrl, authController.RefreshToken).Methods(constant.PostMethod)
	authRouter.HandleFunc(verifyUserUrl, authController.VerifyUser).Methods(constant.GetMethod)

	// Subrouter cho /api/v1
	baseRouter := router.PathPrefix(ApiV1Prefix).Subrouter()
	baseRouter.Use(
		middleware.NewJwtMiddleware(jwt).Do,
	)

	// Subrouter cho /api/v1/users
	userRouter := baseRouter.PathPrefix(UserControllerPrefix).Subrouter()
	userController := user_controller.NewUserController(sqlDb, redis)
	userRouter.HandleFunc(GetUser, userController.GetUser).Methods(constant.GetMethod)
	userRouter.HandleFunc(GetUsers, userController.GetUsers).Methods(constant.GetMethod)
	userRouter.HandleFunc(DeleteUser, userController.DeleteUser).Methods(constant.DeleteMethod)
	userRouter.HandleFunc(UpdateUser, userController.UpdateUser).Methods(constant.PutMethod)
	return router
}
