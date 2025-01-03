package database

import "fmt"

func GetSecret(sessionId string) string {
	RefreshServices()
	result, err := AccountService.GetSession(sessionId)
	if err != nil {
    return fmt.Sprintf("Error: %v", err)
	}
	return result.Secret
}
