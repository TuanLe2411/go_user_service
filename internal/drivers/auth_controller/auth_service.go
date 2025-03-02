package auth_controller

import (
	"encoding/json"
	"go-service-demo/internal/model"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/messaging_system"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	rabbitMq *messaging_system.RabbitMQ
}

func (a *AuthService) CreateVerifyRequest(user model.User) {
	id, _ := uuid.NewUUID()
	verifyRequest := model.UserAccountAction{
		Username:  user.Username,
		Action:    constant.UserVerifyAction,
		CreatedAt: time.Now(),
		RequestID: id.String(),
	}
	rawData, err := json.Marshal(verifyRequest)
	if err != nil {
		return
	}
	a.rabbitMq.PublishWithCtx(rawData)
}
