package repositories

import (
	"go-service-demo/internal/model"
)

type IUserRepo interface {
	Save(model.User) (model.User, error)
	DeleteById(id interface{}) error
	FindById(int) (model.User, error)
	FindAll() ([]model.User, error)
}
