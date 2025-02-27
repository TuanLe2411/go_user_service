package redis

func GetUserKey(username string) string {
	return "user:" + username
}
