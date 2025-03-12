package auth_controller

import (
	"encoding/json"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/database"
	"go-service-demo/pkg/database/mysql/repository_impl"
	"go-service-demo/pkg/messaging_system"
	"go-service-demo/pkg/object"
	"go-service-demo/pkg/utils"
	"log"
	"net/http"
	"strconv"
)

type AuthController struct {
	*AuthService
}

func NewAuthController(db database.Database, jwt *utils.Jwt, rabbitMq *messaging_system.RabbitMQ) *AuthController {
	return &AuthController{
		AuthService: &AuthService{
			rabbitMq:              rabbitMq,
			userRepo:              repository_impl.NewUserRepo(db),
			userAccountActionRepo: repository_impl.NewUserAccountActionRepoImpl(db),
			jwt:                   jwt,
		},
	}
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginUser object.UserLogin
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		log.Println("Error when parse JSON: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := a.userRepo.FindByUsername(loginUser.Username)
	if err != nil {
		log.Println("Error when find user by username: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
	if !user.IsVerified {
		log.Println("User is not verified")
		utils.SetHttpReponseError(r, utils.ErrUnAuthorized)
		return
	}

	if !utils.CheckPasswordHash(loginUser.Password, user.Password) {
		log.Println("Username or password is incorrect")
		utils.SetHttpReponseError(r, utils.ErrUnAuthorized)
		return
	}

	token, err := a.jwt.GenerateAccessToken(strconv.Itoa(user.Id))
	if err != nil {
		log.Println("Error when generate access token: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	refreshToken, err := a.jwt.GenerateRefreshToken(strconv.Itoa(user.Id))
	if err != nil {
		log.Println("Error when generate refresh token: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(object.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (a *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var registerUser object.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&registerUser)
	if err != nil {
		log.Println("Error when parse JSON: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}
	defer r.Body.Close()

	existedUser, err := a.userRepo.FindByUsername(registerUser.Username)
	if err != nil {
		log.Println("Error when find user by username: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
	if existedUser.IsExisted() {
		log.Println("Username is existed")
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	user, err := registerUser.ToUser()
	if err != nil {
		log.Println("Error when convert register user to user: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	verifyRequest := a.createVerifyRequest(user)
	err = a.userAccountActionRepo.Insert(verifyRequest)
	if err != nil {
		log.Println("Error when insert verify request: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	err = a.userRepo.Insert(user)
	if err != nil {
		log.Println("Error when insert user: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
	user.Password = ""

	go a.sendToMessagingSystem(verifyRequest)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (a *AuthController) Logout(w http.ResponseWriter, r *http.Request) {}

func (a *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshToken object.RefreshToken
	err := json.NewDecoder(r.Body).Decode(&refreshToken)
	if err != nil {
		log.Println("Error when parse JSON: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}
	defer r.Body.Close()

	isValid, claims := a.jwt.ValidateRefreshToken(refreshToken.RefreshToken)
	if !isValid {
		log.Println("Invalid refresh token")
		utils.SetHttpReponseError(r, utils.ErrUnAuthorized)
		return
	}

	token, err := a.jwt.GenerateAccessToken(claims.UserId)
	if err != nil {
		log.Println("Error when generate access token: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	w.Write([]byte(token))
}

func (a *AuthController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if len(token) == 0 {
		log.Println("Token is empty")
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}
	userAccountAction, err := a.userAccountActionRepo.FindByRequestId(token)
	if err != nil {
		log.Println("Error when find user account action by request id: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
	if !userAccountAction.IsExisted() {
		log.Println("User account action is not existed")
		utils.SetHttpReponseError(r, utils.ErrNotFound)
		return
	}
	if userAccountAction.Action != constant.UserVerifyAction {
		log.Println("Invalid action")
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	user, err := a.userRepo.FindByUsername(userAccountAction.Username)
	if err != nil {
		log.Println("Error when find user by username: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
	if user.IsVerified {
		log.Println("User is already verified")
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}

	user.IsVerified = true
	err = a.userRepo.VerifyUserByUsername(user.Username)
	if err != nil {
		log.Println("Error when update user: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
	w.Write([]byte("User is verified successfully"))
}
