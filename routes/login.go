package routes

import (
	"cyansky/appwrite"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var reqBody map[string]any
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request body: %s", err),
		})
		return
	}

	result, err := CreateSession(reqBody["email"].(string), reqBody["password"].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Error creating session: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func CreateSession(email, password string) (map[string]interface{}, error) {
	database.RefreshServices()

	sessionsResult, err := database.AccountService.CreateEmailPasswordSession(
		email,
		password,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	userData := database.GetUserData(sessionsResult.UserId)
	userData["token"] = sessionsResult.Secret

	return userData, nil
}
