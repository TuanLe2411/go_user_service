package object

import (
	"go-service-demo/internal/model"
	"go-service-demo/pkg/utils"
)

type RegisterUser struct {
	Username string     `json:"username"`
	Password string     `json:"password"`
	Email    string     `json:"email"`
	UserInfo CreateUser `json:"userInfo"`
}

func (r RegisterUser) ToUser() (model.User, error) {
	hashPassword, err := utils.HashPassword(r.Password)
	if err != nil {
		return model.User{}, err
	}
	return model.User{
		Username:    r.Username,
		Password:    hashPassword,
		Name:        r.UserInfo.Name,
		Age:         r.UserInfo.Age,
		DateOfBirth: r.UserInfo.DateOfBirth,
		Email:       r.Email,
	}, nil
}
