package repositories

import (
	"go-service-demo/internal/model"
)

type IUserRepo interface {
	Save(model.User) (model.User, error)
	// Delete(model.User) error
	// Insert(model.User) (model.User, error)
	FindById(int) (model.User, error)
	FindAll() ([]model.User, error)
}
