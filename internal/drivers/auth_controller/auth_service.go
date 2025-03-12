package auth_controller

import (
	"encoding/json"
	"go-service-demo/internal/model"
	"go-service-demo/internal/repositories"
	"go-service-demo/pkg/constant"
	"go-service-demo/pkg/database/redis"
	"go-service-demo/pkg/messaging_system"
	"go-service-demo/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type AuthService struct {
	rabbitMq              *messaging_system.RabbitMQ
	userRepo              repositories.UserRepo
	userAccountActionRepo repositories.UserAccountActionRepo
	jwt                   *utils.Jwt
	redis                 *redis.RedisDatabase
}

func (a *AuthService) createVerifyRequest(user model.User) model.UserAccountAction {
	id, _ := uuid.NewUUID()
	return model.UserAccountAction{
		Username:  user.Username,
		Action:    constant.UserVerifyAction,
		CreatedAt: time.Now(),
		RequestID: id.String(),
		Email:     user.Email,
	}
}

func (a *AuthService) sendToMessagingSystem(verifyRequest model.UserAccountAction) {
	rawData, _ := json.Marshal(verifyRequest)
	go a.rabbitMq.Publish(rawData)
}
