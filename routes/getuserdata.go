package routes

import (
	database "cyansky/appwrite"
	"fmt"
	"net/http"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/id"
	"github.com/gin-gonic/gin"
)

func GetUserData(c *gin.Context) {
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
	)

	accService := appwrite.NewAccount(userClient)

	response, err := accService.CreateSession(
		id.Unique(),
		reqBody["token"].(string),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
      "error": fmt.Sprintf("Error creating session: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
