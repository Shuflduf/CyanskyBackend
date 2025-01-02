package routes

import (
	"cyansky/appwrite"
	"fmt"
	"net/http"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/query"
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

	userClient := appwrite.NewClient(
		appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
		appwrite.WithProject(database.ProjectId),
		// appwrite.WithKey(os.Getenv("ADMIN_SECRET")),
	)

	accountService := appwrite.NewAccount(userClient)

	sessionsResult, err := accountService.CreateEmailPasswordSession(
		reqBody["email"].(string),
		reqBody["password"].(string),
	)

	if err != nil {
    // TODO: fix all the http statuses to be more descriptive
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Couldn't create session: " + fmt.Sprintf("%v", err),
		})
		return
	}

  userData := GetUserData(sessionsResult.UserId)
  username := userData["username"].(string)

	c.JSON(http.StatusOK, gin.H{
    "message": "Logged in as " + username,
	})
}

func GetUserData(userId string) map[string]interface{} {
	documentList, err := database.DatabaseService.ListDocuments(
		"cyansky-main",
		"user-data",
    database.DatabaseService.WithListDocumentsQueries([]string{query.Equal("auth-id", userId)}),
	)

	if err != nil {
    fmt.Printf("DB: %v", err)
		return nil
	}

	var info map[string]interface{}
	err = documentList.Decode(&info)
  if err != nil {
    fmt.Printf("Decode: %v", err)
    return nil
  }
	return info["documents"].([]interface{})[0].(map[string]interface{})
}
