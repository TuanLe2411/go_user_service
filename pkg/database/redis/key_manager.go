package redis

func GetUserKey(userId string) string {
	return "user:" + userId
}
