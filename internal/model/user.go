package model

type User struct {
	Name        string `json:"name,omitempty"`
	Age         int    `json:"age,omitempty"`
	DateOfBirth string `json:"dateOfBirth,omitempty"`
	Id          int    `json:"id,omitempty"`
}
