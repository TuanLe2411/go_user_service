package repository_impl

import (
	"fmt"
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
	query := "SELECT id, age, name, date_of_birth, username, email FROM user where is_verified = true"
	rows, cancel, err := u.db.QueryRows(query)
	defer cancel()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Age, &user.Name, &user.DateOfBirth, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserRepoImpl) Insert(user model.User) error {
	query := "INSERT INTO user (age, name, date_of_birth, password, username, email) VALUES (?, ?, ?, ?, ?, ?)"
	_, cancel, err := u.db.Exec(query, user.Age, user.Name, user.DateOfBirth, user.Password, user.Username, user.Email)
	defer cancel()
	return err
}

func (u *UserRepoImpl) DeleteByUsername(username string) error {
	query := "DELETE FROM user WHERE username = ?"
	_, cancel, err := u.db.Exec(query, username)
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepoImpl) FindByUsername(username string) (model.User, error) {
	query := "SELECT id, age, name, date_of_birth, username, is_verified, email FROM user WHERE username = ?"
	row, cancel, err := u.db.QueryRow(query, username)
	defer cancel()
	if err != nil {
		return model.User{}, err
	}
	var user model.User
	if err := row.Scan(&user.Id, &user.Age, &user.Name, &user.DateOfBirth, &user.Username, &user.IsVerified, &user.Email); err != nil {
		return model.User{}, nil
	}
	return user, nil

}

func (u *UserRepoImpl) FindPasswordByUsername(username string) (string, error) {
	query := "select user.password from user where username = ?"
	row, cancel, err := u.db.QueryRow(query, username)
	defer cancel()
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var password string
	if err := row.Scan(&password); err != nil {
		return "", err
	}
	return password, nil
}

func (u *UserRepoImpl) UpdateByUsername(user model.User) error {
	query := "UPDATE user SET age = ?, name = ?, date_of_birth = ?, email = ? WHERE username = ?"
	_, cancel, err := u.db.Exec(query, user.Age, user.Name, user.DateOfBirth, user.Email, user.Username)
	defer cancel()
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepoImpl) VerifyUserByUsername(username string) error {
	query := "UPDATE user SET is_verified = true WHERE username = ?"
	_, cancel, err := u.db.Exec(query, username)
	defer cancel()
	return err
}
