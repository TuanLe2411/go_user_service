package user_controller

import (
	"encoding/json"
	"go-service-demo/internal/repositories"
	"go-service-demo/pkg/database"
	"go-service-demo/pkg/database/mysql/repository_impl"
	"go-service-demo/pkg/database/redis"
	"go-service-demo/pkg/object"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserController struct {
	userRepo repositories.IUserRepo
	redis    *redis.RedisDatabase
}

func NewUserController(db database.IDatabase, redis *redis.RedisDatabase) *UserController {
	return &UserController{
		userRepo: repository_impl.NewUserRepo(db),
		redis:    redis,
	}
}

func (u *UserController) FindUserById(w http.ResponseWriter, r *http.Request) {
	// Get user id from path
	rawId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(rawId)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	userInRedis, err := u.redis.Get(redis.GetUserKey(rawId))
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(userInRedis))
		return
	}

	user, err := u.userRepo.FindById(id)
	if err != nil {
		log.Println("Error when find user by id: " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Save user info to redis
	u.redis.SaveUserToRedis(redis.GetUserKey(rawId), user)

	// Return user info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *UserController) FindAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := u.userRepo.FindAll()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return users info
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (u *UserController) SaveUser(w http.ResponseWriter, r *http.Request) {
	var createUser object.CreateUser
	err := json.NewDecoder(r.Body).Decode(&createUser)
	if err != nil {
		http.Error(w, "Can not parse JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := createUser.ToUser()
	user, err = u.userRepo.Save(user)
	if err != nil {
		log.Println("Error when save user: " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get user id from path
	rawId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(rawId)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var updateUser object.UpdateUser
	err = json.NewDecoder(r.Body).Decode(&updateUser)
	if err != nil {
		http.Error(w, "Can not parse JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := updateUser.ToUser()
	user.Id = id
	user, err = u.userRepo.Save(user)
	if err != nil {
		log.Println("Error when save user: " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Save user info to redis
	u.redis.SaveUserToRedis(redis.GetUserKey(rawId), user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get user id from path
	rawId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(rawId)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	err = u.userRepo.DeleteById(id)
	if err != nil {
		log.Println("Error when delete user: " + err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Delete user info in redis
	err = u.redis.Del(redis.GetUserKey(rawId))
	if err != nil {
		log.Println("Error when delete user in redis: " + err.Error())
	}

	w.WriteHeader(http.StatusNoContent)
}
