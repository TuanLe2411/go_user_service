package constant

const (
	GetMethod    string = "GET"
	PostMethod   string = "POST"
	DeleteMethod string = "DELETE"
	PutMethod    string = "PUT"
)

const UserContextKey contextKey = "user"
const AppErrorContextKey contextKey = "appError"
const (
	UserVerifyAction UserAction = "user_verify"
)

// RabbitMQ
const UserActionQueueName string = "user_action"
const UserActionExchangeName string = "user_action_exchange"
const UserActionRoutingKey string = "user_action_routing_key"
