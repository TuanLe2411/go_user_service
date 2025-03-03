package model

type User struct {
	Name        string `json:"name,omitempty"`
	Age         int    `json:"age,omitempty"`
	DateOfBirth string `json:"dateOfBirth,omitempty"`
	Id          int    `json:"id,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	IsVerified  bool   `json:"isVerified,omitempty"`
}

func (u User) IsExisted() bool {
	return len(u.Username) > 0
}
