package routes

import (
	"cyansky/appwrite"
	"fmt"
	"net/http"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/gin-gonic/gin"
)

func SetUserData(c *gin.Context) {
	var reqBody map[string]any
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	description, descriptionOk := reqBody["description"].(string)
	displayName, displayNameOk := reqBody["display_name"].(string)
	updateData := map[string]interface{}{}
	if descriptionOk {
		updateData["description"] = description
	}
	if displayNameOk {
		updateData["display-name"] = displayName
	}

	authorClient := appwrite.NewClient(
		appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
		appwrite.WithProject(database.ProjectId),
		appwrite.WithSession(reqBody["token"].(string)),
	)
	service := appwrite.NewAccount(authorClient)
	accountData, accErr := service.Get()
	if accErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to get account data: %s", accErr),
		})
		return
	}
	fmt.Println(accountData)
	userData := database.GetUserData(accountData.Id)
	documentId := userData["$id"].(string)

	database.RefreshServices()
	_, err := database.DatabaseService.UpdateDocument(
		"cyansky-main",
		"user-data",
		documentId,
		database.DatabaseService.WithUpdateDocumentData(updateData),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to update user data: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User data updated",
	})
}
