package auth_controller

import (
	"encoding/json"
	"go-service-demo/internal/repositories"
	"go-service-demo/pkg/database"
	"go-service-demo/pkg/database/mysql/repository_impl"
	"go-service-demo/pkg/object"
	"go-service-demo/pkg/utils"
	"net/http"
)

type AuthController struct {
	userRepo repositories.IUserRepo
	jwt      *utils.Jwt
}

func NewAuthController(db database.IDatabase, jwt *utils.Jwt) *AuthController {
	return &AuthController{
		userRepo: repository_impl.NewUserRepo(db),
		jwt:      jwt,
	}
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginUser object.UserLogin
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		http.Error(w, "Can not parse JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	savedPassword, err := a.userRepo.FindPasswordByUsername(loginUser.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !utils.CheckPasswordHash(loginUser.Password, savedPassword) {
		http.Error(w, "Username or password is incorrect", http.StatusBadRequest)
		return
	}

	token, err := a.jwt.GenerateAccessToken(loginUser.Username)
	if err != nil {
		http.Error(w, "Get error when generate access token", http.StatusInternalServerError)
		return
	}

	refresgToken, err := a.jwt.GenerateRefreshToken(loginUser.Username)
	if err != nil {
		http.Error(w, "Get error when generate refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(object.LoginResponse{
		Token:        token,
		RefreshToken: refresgToken,
	})
}

func (a *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var registerUser object.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&registerUser)
	if err != nil {
		http.Error(w, "Can not parse JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	existedUser, err := a.userRepo.FindByUsername(registerUser.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if existedUser.IsExisted() {
		http.Error(w, "Username is existed", http.StatusBadRequest)
		return
	}

	user, err := registerUser.ToUser()
	if err != nil {
		http.Error(w, "Can not parse JSON", http.StatusBadRequest)
		return
	}
	savedUser, err := a.userRepo.Insert(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(savedUser)
}

func (a *AuthController) Logout(w http.ResponseWriter, r *http.Request) {}

func (a *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshToken object.RefreshToken
	err := json.NewDecoder(r.Body).Decode(&refreshToken)
	if err != nil {
		http.Error(w, "Can not parse JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	isValid, claims := a.jwt.ValidateRefreshToken(refreshToken.RefreshToken)
	if !isValid {
		http.Error(w, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	token, err := a.jwt.GenerateAccessToken(claims.Username)
	if err != nil {
		http.Error(w, "Get error when generate access token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(token))
}
