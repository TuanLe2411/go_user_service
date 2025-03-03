package repository_impl

import (
	"go-service-demo/internal/model"
	"go-service-demo/internal/repositories"
	"go-service-demo/pkg/database"
)

type UserRepoImpl struct {
	db database.Database
}

func NewUserRepo(db database.Database) repositories.UserRepo {
	return &UserRepoImpl{
		db: db,
	}
}

func (u *UserRepoImpl) FindAll() ([]model.User, error) {
	query := "SELECT id, age, name, date_of_birth, username FROM user"
	rows, err := u.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Age, &user.Name, &user.DateOfBirth, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserRepoImpl) Insert(user model.User) error {
	query := "INSERT INTO user (age, name, date_of_birth, password, username) VALUES (?, ?, ?, ?, ?)"
	_, err := u.db.Exec(query, user.Age, user.Name, user.DateOfBirth, user.Password, user.Username)
	return err
}

func (u *UserRepoImpl) DeleteByUsername(username string) error {
	query := "DELETE FROM user WHERE username = ?"
	_, err := u.db.Exec(query, username)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepoImpl) FindByUsername(username string) (model.User, error) {
	query := "SELECT id, age, name, date_of_birth, username, is_verified FROM user WHERE username = ?"
	row, err := u.db.QueryRow(query, username)
	if err != nil {
		return model.User{}, err
	}
	var user model.User
	if err := row.Scan(&user.Id, &user.Age, &user.Name, &user.DateOfBirth, &user.Username, &user.IsVerified); err != nil {
		return model.User{}, err
	}
	return user, nil

}

func (u *UserRepoImpl) FindPasswordByUsername(username string) (string, error) {
	query := "SELECT password FROM user WHERE username = ?"
	row, err := u.db.QueryRow(query, username)
	if err != nil {
		return "", err
	}
	var password string
	if err := row.Scan(&password); err != nil {
		return "", err
	}
	return password, nil
}

func (u *UserRepoImpl) UpdateByUsername(user model.User) error {
	query := "UPDATE user SET age = ?, name = ?, date_of_birth = ?, is_verified = ? WHERE username = ?"
	_, err := u.db.Exec(query, user.Age, user.Name, user.DateOfBirth, user.IsVerified, user.Username)
	if err != nil {
		return err
	}
	return nil
}
