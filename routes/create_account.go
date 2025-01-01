package routes

import (
	"cyansky/appwrite"
	"fmt"
	"net/http"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/gin-gonic/gin"
)

func CreateAccount(c *gin.Context) {
	var reqBody map[string]any
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// register acc into AUTH
	newAccount := appwrite.NewAccount(database.BaseClient)
	result, err := newAccount.Create(
		id.Unique(),
		reqBody["email"].(string),
		reqBody["password"].(string),
		database.AccountManager.WithCreateName(reqBody["name"].(string)),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("%s", err),
		})
		return
	}

	// register acc into DB
	_, _ = database.DatabaseManager.CreateDocument(
		"cyansky-main",
		"user-data",
		id.Unique(),
		map[string]interface{}{
			"username": reqBody["name"].(string),
			"auth-id":  result.Id,
		},
	)

	c.JSON(http.StatusOK, gin.H{
		"message": result,
	})

}
