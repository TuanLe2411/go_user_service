package router

import (
	"go-service-demo/internal/controllers/user_controller"
	"go-service-demo/internal/middleware"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/database"

	"github.com/gorilla/mux"
)

const ApiV1Prefix = "/api/v1"

const UserControllerPrefix = "/users"

const FindUserByIdUrl = "/{id}"
const FindAllUsersUrl = ""
const CreateNewUserUrl = ""

func InitRouter(db database.ISqlDatabase) *mux.Router {
	router := mux.NewRouter()

	// Subrouter cho /users
	userRouter := router.PathPrefix(ApiV1Prefix + UserControllerPrefix).Subrouter()
	userRouter.Use(
		middleware.NewTrackingMiddleware().Do,
		middleware.NewMonitorMiddleware().Do,
	)
	userRouter.HandleFunc(FindUserByIdUrl, user_controller.NewUserController(db).FindUserById).Methods(constant.GetMethod)
	userRouter.HandleFunc(FindAllUsersUrl, user_controller.NewUserController(db).FindAllUsers).Methods(constant.GetMethod)
	userRouter.HandleFunc(CreateNewUserUrl, user_controller.NewUserController(db).SaveUser).Methods(constant.PostMethod)

	return router
}
