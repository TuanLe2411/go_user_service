package object

import "go-service-demo/internal/model"

type CreateUser struct {
	Name        string
	Age         int
	DateOfBirth string
	Email       string
}

func (c CreateUser) ToUser() model.User {
	return model.User{
		Name:        c.Name,
		Age:         c.Age,
		DateOfBirth: c.DateOfBirth,
		Email:       c.Email,
	}
}
