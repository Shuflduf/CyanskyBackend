package database

func GetSecret(sessionId string) string {
  result, err := AccountService.GetSession(sessionId)
  if err != nil {
    return ""
  }
  return result.Secret
}
