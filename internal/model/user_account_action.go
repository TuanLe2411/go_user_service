package model

import (
	"go-service-demo/pkg/constant"
	"time"
)

type UserAccountAction struct {
	ID        int                 `json:"id"`
	Username  string              `json:"username"`
	Action    constant.UserAction `json:"action"`
	CreatedAt time.Time           `json:"createdAt"`
	RequestID string              `json:"requestId"`
}

func (u UserAccountAction) IsExisted() bool {
	return len(u.Username) > 0
}
