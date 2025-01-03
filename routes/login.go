package routes

import (
	"cyansky/appwrite"
	"net/http"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var reqBody map[string]any
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	result, err := CreateSession(reqBody["email"].(string), reqBody["password"].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

  c.JSON(http.StatusOK, result)
}

func CreateSession(email, password string) (map[string]interface{}, error) {
  database.RefreshServices()

	accountService := appwrite.NewAccount(database.AdminClient)

	sessionsResult, err := accountService.CreateEmailPasswordSession(
		email,
		password,
	)

	if err != nil {
		return nil, err
	}

	userData := database.GetUserData(sessionsResult.UserId)
  userData["token"] = sessionsResult.Secret

  return userData, nil
}

