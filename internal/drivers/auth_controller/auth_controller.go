package auth_controller

import (
	"encoding/json"
	"go-service-demo/internal/repositories"
	"go-service-demo/pkg/database"
	"go-service-demo/pkg/database/mysql/repository_impl"
	"go-service-demo/pkg/messaging_system"
	"go-service-demo/pkg/object"
	"go-service-demo/pkg/utils"
	"log"
	"net/http"
)

type AuthController struct {
	userRepo repositories.IUserRepo
	jwt      *utils.Jwt
	*AuthService
}

func NewAuthController(db database.IDatabase, jwt *utils.Jwt, rabbitMq *messaging_system.RabbitMQ) *AuthController {
	return &AuthController{
		userRepo: repository_impl.NewUserRepo(db),
		jwt:      jwt,
		AuthService: &AuthService{
			rabbitMq: rabbitMq,
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

	token, err := a.jwt.GenerateAccessToken(loginUser.Username)
	if err != nil {
		log.Println("Error when generate access token: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	refreshToken, err := a.jwt.GenerateRefreshToken(loginUser.Username)
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
	savedUser, err := a.userRepo.Insert(user)
	if err != nil {
		log.Println("Error when insert user: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	go a.CreateVerifyRequest(savedUser)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(savedUser)
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

	token, err := a.jwt.GenerateAccessToken(claims.Username)
	if err != nil {
		log.Println("Error when generate access token: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	w.Write([]byte(token))
}

func (a *AuthController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Verify user"))
}
