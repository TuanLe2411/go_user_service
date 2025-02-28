package model

import (
	"go-service-demo/pkg/constant"
	"time"
)

type UserAccountAction struct {
	ID        int                 `json:"id"`
	UserID    int                 `json:"userId"`
	Action    constant.UserAction `json:"action"`
	CreatedAt time.Time           `json:"createdAt"`
	RequestID string              `json:"requestId"`
}
