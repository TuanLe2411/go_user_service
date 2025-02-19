package repository_impl

import (
	"go-service-demo/internal/model"
	"go-service-demo/internal/repository"
	"go-service-demo/pkg/database"
)

type UserRepoImpl struct {
	db database.Database
}

func NewUserRepo(db database.Database) repository.UserRepo {
	return &UserRepoImpl{
		db: db,
	}
}

func (u *UserRepoImpl) Save(model.User) (model.User, error) {
	return model.User{}, nil
}

func (u *UserRepoImpl) Delete(model.User) error {
	return nil
}

func (u *UserRepoImpl) Insert(model.User) (model.User, error) {
	return model.User{}, nil
}

func (u *UserRepoImpl) FindOneByName(string) (model.User, error) {
	return model.User{}, nil
}

func (u *UserRepoImpl) FindAll() ([]model.User, error) {
	return []model.User{
		{
			Id:      "1",
			Name:    "Tuan Le",
			Age:     25,
			Address: "Ha noi",
		},
	}, nil
}
