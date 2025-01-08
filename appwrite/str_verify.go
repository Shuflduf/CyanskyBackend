package database

func VerifyUsername(username string) bool {
	for _, char := range username {
		if (char < 'a' || char > 'z') && (char < '0' || char > '9') {
			return false
		}
	}
	return true
}
