package routes

import (
	"cyansky/appwrite"
	"fmt"
	"net/http"

	"github.com/appwrite/sdk-for-go/query"
	"github.com/gin-gonic/gin"
)

func GetUserDataId(c *gin.Context) {
	var reqBody map[string]any
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	username, hasUsername := reqBody["username"].(string)
	userId, hasUserId := reqBody["user_id"].(string)

	database.RefreshServices()

	var info map[string]interface{}
	if !hasUsername && hasUserId {
		var err error
		document, err := database.DatabaseService.GetDocument(
			"cyansky-main",
			"user-data",
			userId,
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Error getting document list: %s", err),
			})
			return
		}
		err = document.Decode(&info)
	} else {
		documentList, err := database.DatabaseService.ListDocuments(
			"cyansky-main",
			"user-data",
			database.DatabaseService.WithListDocumentsQueries([]string{query.Equal("username", username)}),
		)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Error getting document list: %s", err),
			})
			return
		}

    var documentListData map[string]interface{}
		err = documentList.Decode(&documentListData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Error decoding document: %s", err),
			})
			return
		}
    info = documentListData["documents"].([]interface{})[0].(map[string]interface{})
	}

	c.JSON(http.StatusOK, info)
}
