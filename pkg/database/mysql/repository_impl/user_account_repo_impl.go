package repository_impl

import (
	"go-service-demo/internal/model"
	"go-service-demo/internal/repositories"
	"go-service-demo/pkg/database"
)

type UserAccountActionRepoImpl struct {
	db database.Database
}

func NewUserAccountActionRepoImpl(db database.Database) repositories.UserAccountActionRepo {
	return &UserAccountActionRepoImpl{
		db: db,
	}
}

func (u *UserAccountActionRepoImpl) Insert(userAccountAction model.UserAccountAction) error {
	query := "INSERT INTO user_account_action (username, request_id, action, create_at) VALUES (?, ?, ?, ?)"
	_, err := u.db.Exec(query, userAccountAction.Username, userAccountAction.RequestID, userAccountAction.Action, userAccountAction.CreatedAt.Format("2006-01-02 15:04:05"))
	return err
}

func (u *UserAccountActionRepoImpl) FindByRequestId(requestId string) (model.UserAccountAction, error) {
	query := "SELECT id, username, request_id, action, create_at FROM user_account_action WHERE request_id = ?"
	row, err := u.db.QueryRow(query, requestId)
	if err != nil {
		return model.UserAccountAction{}, err
	}
	var userAccountAction model.UserAccountAction
	row.Scan(&userAccountAction.ID, &userAccountAction.Username, &userAccountAction.RequestID, &userAccountAction.Action, &userAccountAction.CreatedAt)
	return userAccountAction, err
}
