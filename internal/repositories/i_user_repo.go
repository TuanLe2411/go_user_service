package repositories

import (
	"go-service-demo/internal/model"
)

type IUserRepo interface {
	Insert(model.User) (model.User, error)
	DeleteById(id interface{}) error
	FindById(int) (model.User, error)
	FindByUsername(string) (model.User, error)
	FindAll() ([]model.User, error)
	FindPasswordByUsername(string) (string, error)
	UpdateByUsername(model.User) error
}
