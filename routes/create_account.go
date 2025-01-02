package routes

import (
	"cyansky/appwrite"
	"fmt"
	"net/http"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/id"
	// "github.com/appwrite/sdk-for-go/role"
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

	newClient := appwrite.NewClient(
		appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
		appwrite.WithProject(database.ProjectId),
	)

	// register acc into AUTH
	newAccount := appwrite.NewAccount(newClient)
	accountResult, accErr := newAccount.Create(
		id.Unique(),
		reqBody["email"].(string),
		reqBody["password"].(string),
		database.AccountService.WithCreateName(reqBody["name"].(string)),
	)

	if accErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to create account: %s", accErr),
		})
		return
	}
	userId := accountResult.Id

	// register acc into DB
	result, err := database.DatabaseService.CreateDocument(
		"cyansky-main",
		"user-data",
		id.Unique(),
		map[string]interface{}{
			"username": reqBody["name"].(string),
			"auth-id":  userId,
		},
		database.DatabaseService.WithCreateDocumentPermissions([]string{
			"write(\"user:" + userId + "\")",
			"read(\"user:" + userId + "\")",
		}),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to register account into database: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
