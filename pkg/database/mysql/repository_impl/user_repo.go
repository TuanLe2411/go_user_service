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
	query := "SELECT id, age, name, date_of_birth FROM user WHERE id = %d"
	rows, err := u.db.Query(fmt.Sprintf(query, id))
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Age, &user.Name, &user.DateOfBirth); err != nil {
			log.Fatal(err)
		}
		return user, nil
	}
	return model.User{}, nil
}

func (u *UserRepo) FindAll() ([]model.User, error) {
	query := "SELECT id, age, name, date_of_birth FROM user"
	rows, err := u.db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.Id, &user.Age, &user.Name, &user.DateOfBirth); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserRepo) Save(user model.User) (model.User, error) {
	query := "INSERT INTO user (age, name, date_of_birth) VALUES (%d, '%s', '%s')"
	_, err := u.db.Query(fmt.Sprintf(query, user.Age, user.Name, user.DateOfBirth))
	if err != nil {
		log.Fatal(err)
	}
	return user, nil
}
