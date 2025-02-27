package repository_impl

import (
	"fmt"
	"go-service-demo/internal/model"
	"go-service-demo/internal/repositories"
	"go-service-demo/pkg/database"
	"log"
)

type UserRepo struct {
	db database.IDatabase
}

func NewUserRepo(db database.IDatabase) repositories.IUserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) FindById(id int) (model.User, error) {
	query := "SELECT id, age, name, date_of_birth, username FROM user WHERE id = %d"
	rows, err := u.db.Query(fmt.Sprintf(query, id))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Age, &user.Name, &user.DateOfBirth, &user.Username); err != nil {
			return model.User{}, err
		}
		return user, nil
	}
	return model.User{}, nil
}

func (u *UserRepo) FindAll() ([]model.User, error) {
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

func (u *UserRepo) Insert(user model.User) (model.User, error) {
	query := "INSERT INTO user (age, name, date_of_birth, password, username) VALUES (%d, '%s', '%s', '%s', '%s')"
	_, err := u.db.Query(fmt.Sprintf(query, user.Age, user.Name, user.DateOfBirth, user.Password, user.Username))
	if err != nil {
		return model.User{}, err
	}
	user.Password = ""
	return user, nil
}

func (u *UserRepo) DeleteById(id interface{}) error {
	query := "DELETE FROM user WHERE id = %d"
	_, err := u.db.Query(fmt.Sprintf(query, id))
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) FindByUsername(username string) (model.User, error) {
	query := "SELECT id, age, name, date_of_birth, username FROM user WHERE username = '%s'"
	rows, err := u.db.Query(fmt.Sprintf(query, username))
	if err != nil {
		return model.User{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Age, &user.Name, &user.DateOfBirth, &user.Username); err != nil {
			return model.User{}, err
		}
		return user, nil
	}
	return model.User{}, nil
}

func (u *UserRepo) FindPasswordByUsername(username string) (string, error) {
	query := "SELECT password FROM user WHERE username = '%s'"
	rows, err := u.db.Query(fmt.Sprintf(query, username))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		var password string
		if err := rows.Scan(&password); err != nil {
			return "", err
		}
		return password, nil
	}
	return "", nil
}

func (u *UserRepo) UpdateByUsername(user model.User) error {
	query := "UPDATE user SET age = %d, name = '%s', date_of_birth = '%s' WHERE username = '%s'"
	_, err := u.db.Query(fmt.Sprintf(query, user.Age, user.Name, user.DateOfBirth, user.Username))
	if err != nil {
		return err
	}
	return nil
}
