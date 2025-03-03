package repositories

import "go-service-demo/internal/model"

type UserAccountActionRepo interface {
	Insert(model.UserAccountAction) error
	FindByRequestId(string) (model.UserAccountAction, error)
}
