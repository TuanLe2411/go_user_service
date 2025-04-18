package user_controller

import (
	"encoding/json"
	"errors"
	"go-service-demo/internal/repositories"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/database"
	"go-service-demo/pkg/database/mysql/repository_impl"
	"go-service-demo/pkg/database/redis"
	"go-service-demo/pkg/object"
	"go-service-demo/pkg/utils"
	"net/http"

	"github.com/rs/zerolog/log"
)

type UserController struct {
	userRepo repositories.UserRepo
	redis    *redis.RedisDatabase
}

func NewUserController(
	db database.Database,
	redis *redis.RedisDatabase,
) *UserController {
	return &UserController{
		userRepo: repository_impl.NewUserRepo(db),
		redis:    redis,
	}
}

func (u *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	username := r.Header.Get(constant.UsernameHeaderKey)
	userInRedis, err := u.redis.Get(redis.GetUserKey(username))
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(userInRedis))
		return
	}

	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when find user by username: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	// Save user info to redis
	u.redis.SaveUserToRedis(redis.GetUserKey(username), user)

	// Return user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	users, err := u.userRepo.FindAll()
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when find all users: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}
	if len(users) == 0 {
		msg := "no user found"
		log.Error().
			Str("trackingId", trackingId).
			Str("error", msg).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrNotFound, errors.New(msg))
		return
	}

	// Return users info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	username := r.Header.Get(constant.UsernameHeaderKey)

	var updateUser object.UpdateUser
	err := json.NewDecoder(r.Body).Decode(&updateUser)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when parse JSON: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrBadRequest, err)
		return
	}
	defer r.Body.Close()

	user := updateUser.ToUser()
	user.Username = username
	err = u.userRepo.UpdateByUsername(user)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when update user: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	// Save user info to redis
	u.redis.SaveUserToRedis(redis.GetUserKey(username), user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	trackingId := r.Context().Value(constant.TrackingIdContextKey).(string)
	username := r.Header.Get(constant.UsernameHeaderKey)

	err := u.userRepo.DeleteByUsername(username)
	if err != nil {
		log.Error().
			Str("trackingId", trackingId).
			Str("error", "Error when delete user: "+err.Error()).
			Msg("")
		utils.SetHttpReponseError(r, utils.ErrServerError, err)
		return
	}

	go func() {
		err = u.redis.Del(redis.GetUserKey(username))
		if err != nil {
			log.Error().
				Str("trackingId", trackingId).
				Str("error", "Error when delete user in redis: "+err.Error()).
				Msg("")
		}
	}()

	w.WriteHeader(http.StatusNoContent)
}
