package user_controller

import (
	"encoding/json"
	"go-service-demo/internal/repositories"
	"go-service-demo/pkg/database"
	"go-service-demo/pkg/database/mysql/repository_impl"
	"go-service-demo/pkg/object"
	"go-service-demo/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserController struct {
	userRepo repositories.IUserRepo
}

func NewUserController(db database.ISqlDatabase) *UserController {
	return &UserController{
		userRepo: repository_impl.NewUserRepo(db),
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

	user, err := u.userRepo.FindById(id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

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
		http.Error(w, "Không thể parse JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user := utils.ToUser(createUser)
	user, err = u.userRepo.Save(user)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
