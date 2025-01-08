package routes

import (
	database "cyansky/appwrite"
	"fmt"
	"net/http"

	"github.com/appwrite/sdk-for-go/query"
	"github.com/gin-gonic/gin"
)

func GetPosts(c *gin.Context) {
	var reqBody map[string]any
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	database.RefreshServices()
	queries := []string{
		query.OrderDesc("$createdAt"),
		query.Limit(10),
	}

	pointerId, isContinue := reqBody["last_post_id"].(string)
	if isContinue && pointerId != "" {
		queries = append(queries, query.CursorAfter(pointerId))
	}

  userId, isUser := reqBody["user_id"].(string)
  if isUser && userId != "" {
    // queries = append(queries, query.Equal()("$userId", "=", userId))
  }

	documentList, err := database.DatabaseService.ListDocuments(
		"cyansky-main",
		"posts",
		database.DatabaseService.WithListDocumentsQueries(queries),
	)

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": fmt.Sprintf("Internal server error: %v", err),
    })
    return
  }

	var decodedDocuments map[string]interface{}
	err = documentList.Decode(&decodedDocuments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to decode documents: %s", err),
		})
	}
  c.JSON(http.StatusOK, decodedDocuments)
}
