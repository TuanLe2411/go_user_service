package user_controller

import (
	"encoding/json"
	"go-service-demo/internal/repositories"
	"go-service-demo/pkg/database"
	"go-service-demo/pkg/database/mysql/repository_impl"
	"go-service-demo/pkg/database/redis"
	"go-service-demo/pkg/object"
	"go-service-demo/pkg/utils"
	"log"
	"net/http"
)

type UserController struct {
	userRepo repositories.UserRepo
	redis    *redis.RedisDatabase
}

func NewUserController(db database.Database, redis *redis.RedisDatabase) *UserController {
	return &UserController{
		userRepo: repository_impl.NewUserRepo(db),
		redis:    redis,
	}
}

func (u *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("user_id")
	userInRedis, err := u.redis.Get(redis.GetUserKey(username))
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(userInRedis))
		return
	}

	user, err := u.userRepo.FindByUsername(username)
	if err != nil {
		log.Println("Error when find user by username: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	// Save user info to redis
	u.redis.SaveUserToRedis(redis.GetUserKey(username), user)

	// Return user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.userRepo.FindAll()
	if err != nil {
		log.Println("Error when find all users: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}
	if len(users) == 0 {
		log.Println("No user found")
		utils.SetHttpReponseError(r, utils.ErrNotFound)
		return
	}

	// Return users info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("user_id")

	var updateUser object.UpdateUser
	err := json.NewDecoder(r.Body).Decode(&updateUser)
	if err != nil {
		log.Println("Error when parse JSON: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrBadRequest)
		return
	}
	defer r.Body.Close()

	user := updateUser.ToUser()
	user.Username = username
	err = u.userRepo.UpdateByUsername(user)
	if err != nil {
		log.Println("Error when update user: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	// Save user info to redis
	u.redis.SaveUserToRedis(redis.GetUserKey(username), user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("user_id")

	err := u.userRepo.DeleteByUsername(username)
	if err != nil {
		log.Println("Error when delete user: " + err.Error())
		utils.SetHttpReponseError(r, utils.ErrServerError)
		return
	}

	// Delete user info in redis
	err = u.redis.Del(redis.GetUserKey(username))
	if err != nil {
		log.Println("Error when delete user in redis: " + err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
}
