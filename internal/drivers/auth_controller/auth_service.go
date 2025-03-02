package auth_controller

import (
	"encoding/json"
	"go-service-demo/internal/model"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/messaging_system"
	"time"
)

type AuthService struct {
	rabbitMq *messaging_system.RabbitMQ
}

func (a *AuthService) CreateVerifyRequest(user model.User) {
	verifyRequest := model.UserAccountAction{
		UserID:    user.Id,
		Action:    constant.UserVerifyAction,
		CreatedAt: time.Now(),
		RequestID: "",
	}
	rawData, err := json.Marshal(verifyRequest)
	if err != nil {
		return
	}
	a.rabbitMq.PublishWithCtx(rawData)
}
