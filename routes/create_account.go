package routes

import (
	"cyansky/appwrite"
	"fmt"
	"net/http"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/appwrite/sdk-for-go/query"

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

	username := reqBody["username"].(string)
	email := reqBody["email"].(string)
	password := reqBody["password"].(string)
	displayName := reqBody["display_name"].(string)

	// check if username is valid
	if !database.VerifyUsername(username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid username, only lowercase letters and numbers are allowed",
		})
		return
	}

	// check if username already exists
	database.RefreshServices()
	usersWithSameName, err := database.DatabaseService.ListDocuments(
		"cyansky-main",
		"user-data",
		database.DatabaseService.WithListDocumentsQueries([]string{
			query.Equal("username", username),
		}),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to check if username exists: %s", err),
		})
		return
	}
	if usersWithSameName.Total > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username already exists",
		})
		return
	}

	database.RefreshServices()
	newClient := appwrite.NewClient(
		appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
		appwrite.WithProject(database.ProjectId),
	)

	// register acc into AUTH
	newAccount := appwrite.NewAccount(newClient)
	accountResult, accErr := newAccount.Create(
		id.Unique(),
		email,
		password,
		database.AccountService.WithCreateName(username),
	)

	if accErr != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": fmt.Sprintf("Failed to create account: %s", accErr),
		})
		return
	}
	userId := accountResult.Id

	// register acc into DB
	_, err = database.DatabaseService.CreateDocument(
		"cyansky-main",
		"user-data",
		id.Unique(),
		map[string]interface{}{
			"username": username,
			"display-name": displayName,
			"auth-id":  userId,
		},
		database.DatabaseService.WithCreateDocumentPermissions([]string{
			"write(\"user:" + userId + "\")",
			"read(\"user:" + userId + "\")",
		}),
	)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": fmt.Sprintf("Failed to register account into database: %s", err),
		})
		return
	}

	// create session
	sessionsResult, err := CreateSession(email, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, sessionsResult)
}
