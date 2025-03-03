package repositories

import (
	"go-service-demo/internal/model"
)

type UserRepo interface {
	Insert(model.User) error
	DeleteByUsername(string) error
	FindByUsername(string) (model.User, error)
	FindAll() ([]model.User, error)
	FindPasswordByUsername(string) (string, error)
	UpdateByUsername(model.User) error
	VerifyUserByUsername(string) error
}
