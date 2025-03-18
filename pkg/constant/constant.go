package constant

const (
	GetMethod    string = "GET"
	PostMethod   string = "POST"
	DeleteMethod string = "DELETE"
	PutMethod    string = "PUT"
)

const UserContextKey contextKey = "user"
const AppErrorContextKey contextKey = "appError"
const TrackingIdContextKey contextKey = "trackingId"

const (
	UserVerifyAction UserAction = "user_verify"
)

const UsernameHeaderKey string = "username"
const UserIdHeaderKey string = "user_id"
