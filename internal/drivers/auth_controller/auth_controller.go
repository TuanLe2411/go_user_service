package auth_controller

import (
	"encoding/json"
	"errors"
	"go-service-demo/internal/model"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/database"
	"go-service-demo/pkg/database/mysql/repository_impl"
	"go-service-demo/pkg/database/redis"
	"go-service-demo/pkg/messaging_system"
	"go-service-demo/pkg/object"
	"go-service-demo/pkg/utils"
	"net/http"

	"github.com/rs/zerolog/log"
)

type AuthController struct {
	*AuthService
}

func NewAuthController(
	db database.Database,
	jwt *utils.Jwt,
	rabbitMq *messaging_system.RabbitMQ,
	redis *redis.RedisDatabase,
) *AuthController {
	return &AuthController{
		AuthService: &AuthService{
			rabbitMq:              rabbitMq,
			userRepo:              repository_impl.NewUserRepo(db),
			userAccountActionRepo: repository_impl.NewUserAccountActionRepoImpl(db),
			jwt:                   jwt,
			redis:                 redis,
		},
	}
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	var loginUser object.UserLogin
	err := json.NewDecoder(r.Body).Decode(&loginUser)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when parse UserLogin object: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, err)
		return
	}
	defer r.Body.Close()

	var user model.User
	cachedUser, err := a.redis.Get(redis.GetUserKey(loginUser.Username))
	if err == nil && len(cachedUser) > 0 {
		_ = json.Unmarshal([]byte(cachedUser), &user)
	} else {
		user, err = a.userRepo.FindByUsername(loginUser.Username)
	}
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when find user by username: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}
	if !user.IsVerified {
		msg := "user is not verified"
		log.Info().
			Str("trackingId", trackingId).
			Str("error", msg).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrUnAuthorized, errors.New(msg))
		return
	}

	password, err := a.userRepo.FindPasswordByUsername(loginUser.Username)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when find user password: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	if !utils.CheckPasswordHash(loginUser.Password, password) {
		msg := "username or password is incorrect"
		log.Info().
			Str("trackingId", trackingId).
			Str("error", msg).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrUnAuthorized, errors.New(msg))
		return
	}

	token, err := a.jwt.GenerateAccessToken(user)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when generate access token: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	refreshToken, err := a.jwt.GenerateRefreshToken(user)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when generate refresh token: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(object.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (a *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	var registerUser object.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&registerUser)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when parse RegisterUser object: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, err)
		return
	}
	defer r.Body.Close()

	existedUser, err := a.userRepo.FindByUsername(registerUser.Username)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when find user by username: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}
	if existedUser.IsExisted() {
		log.Info().
			Str("trackingId", trackingId).
			Str("error", "username is already existed").
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, err)
		return
	}

	user, err := registerUser.ToUser()
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when convert register user to user: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, err)
		return
	}

	verifyRequest := a.createVerifyRequest(user)
	err = a.userAccountActionRepo.Insert(verifyRequest)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when insert verify request: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	err = a.userRepo.Insert(user)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when insert user: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}
	user.Password = ""

	go a.redis.SaveUserToRedis(redis.GetUserKey(user.Username), user)
	go a.sendToMessagingSystem(verifyRequest)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (a *AuthController) Logout(w http.ResponseWriter, r *http.Request) {}

func (a *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	var refreshToken object.RefreshToken
	err := json.NewDecoder(r.Body).Decode(&refreshToken)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when parse JSON: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, err)
		return
	}
	defer r.Body.Close()

	isValid, claims := a.jwt.ValidateRefreshToken(refreshToken.RefreshToken)
	if !isValid {
		log.Info().
			Str("trackingId", trackingId).
			Str("error", "Invalid refresh token").
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrUnAuthorized, err)
		return
	}

	token, err := a.jwt.GenerateAccessToken(model.User{
		Id:       claims.UserId,
		Username: claims.Username,
	})
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when generate access token: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	w.Write([]byte(token))
}

func (a *AuthController) VerifyUser(w http.ResponseWriter, r *http.Request) {
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	token := r.URL.Query().Get("token")
	if len(token) == 0 {
		msg := "token is empty"
		log.Info().
			Str("trackingId", trackingId).
			Str("error", msg).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, errors.New(msg))
		return
	}
	userAccountAction, err := a.userAccountActionRepo.FindByRequestId(token)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when find user account action by request id: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}
	if !userAccountAction.IsExisted() {
		log.Info().
			Str("trackingId", trackingId).
			Str("error", "user account action is not existed").
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrNotFound, err)
		return
	}
	if userAccountAction.Action != constant.UserVerifyAction {
		msg := "invalid action"
		log.Info().
			Str("trackingId", trackingId).
			Str("error", msg).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, errors.New(msg))
		return
	}

	var user model.User
	cachedUser, err := a.redis.Get(redis.GetUserKey(userAccountAction.Username))
	if err == nil && len(cachedUser) > 0 {
		_ = json.Unmarshal([]byte(cachedUser), &user)
	} else {
		user, err = a.userRepo.FindByUsername(userAccountAction.Username)
	}

	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when find user by username: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	if user.IsVerified {
		w.Write([]byte("User is verified successfully"))
		return
	}

	user.IsVerified = true
	err = a.userRepo.VerifyUserByUsername(user.Username)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when update user: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	err = a.redis.SaveUserToRedis(redis.GetUserKey(user.Username), user)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when cached user: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	w.Write([]byte("User is verified successfully"))
}
