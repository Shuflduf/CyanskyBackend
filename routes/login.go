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

	client := appwrite.NewClient(
		appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
		appwrite.WithProject(database.ProjectId),
	)

	accountService := appwrite.NewAccount(client)

	result, err := accountService.CreateEmailPasswordSession(
		reqBody["email"].(string),
		reqBody["password"].(string),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}

	searchQuery := []string{query.Equal("auth-id", []interface{}{result.UserId})}
	documentList, err := database.DatabaseService.ListDocuments(
		"cyansky-main",
		"user-data",
		database.DatabaseService.WithListDocumentsQueries(searchQuery),
	)

	if err != nil {
		fmt.Println(err)
	}

	var info map[string]any
  err = documentList.Documents[0].Decode(&info)
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(documentList.Documents[0])
  fmt.Println(info)

	c.JSON(http.StatusOK, gin.H{
		"message": info,
	})
}
