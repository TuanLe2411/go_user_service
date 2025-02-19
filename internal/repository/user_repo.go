package repository

import (
	"go-service-demo/internal/model"
)

type UserRepo interface {
	Save(model.User) (model.User, error)
	Delete(model.User) error
	Insert(model.User) (model.User, error)
	FindOneByName(string) (model.User, error)
	FindAll() ([]model.User, error)
}
