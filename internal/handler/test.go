package handler

import (
	"encoding/json"
	"go-service-demo/internal/repository"
	"go-service-demo/pkg/database"
	"go-service-demo/pkg/database/mock/repository_impl"
	"log"
	"net/http"
)

type GetTestHandler struct {
	userRepo repository.UserRepo
}

func NewGetTestHandler(db database.Database) *GetTestHandler {
	userRepo := repository_impl.NewUserRepo(db)
	return &GetTestHandler{
		userRepo: userRepo,
	}
}

func (t *GetTestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	users, err := t.userRepo.FindAll()
	if err != nil {
		http.Error(w, "Error get users", http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error marshalling users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonData); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
