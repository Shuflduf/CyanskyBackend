package routes

import (
	"cyansky/appwrite"
	"fmt"
	"net/http"

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

	database.RefreshServices()
	documentList, err := database.DatabaseService.GetDocument(
		"cyansky-main",
		"user-data",
		reqBody["user_id"].(string),
	)

	if err != nil {
		fmt.Printf("DB: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Error getting document list: %s", err),
		})
		return
	}

	var info map[string]interface{}
	err = documentList.Decode(&info)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Error decoding document list: %s", err),
		})
		return
	}
	// output := info["documents"].([]interface{})[0].(map[string]interface{})
	c.JSON(http.StatusOK, info)
}
