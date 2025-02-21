package utils

import (
	"go-service-demo/internal/model"
	"go-service-demo/pkg/object"
)

func ToUser(c object.CreateUser) model.User {
	return model.User{
		Name:        c.Name,
		Age:         c.Age,
		DateOfBirth: c.DateOfBirth,
	}
}
