package object

import "go-service-demo/internal/model"

type UpdateUser struct {
	Name        string
	Age         int
	DateOfBirth string
}

func (c UpdateUser) ToUser() model.User {
	return model.User{
		Name:        c.Name,
		Age:         c.Age,
		DateOfBirth: c.DateOfBirth,
	}
}
