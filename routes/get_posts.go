package routes

import (
	database "cyansky/appwrite"
	"fmt"
	"net/http"
	"sort"
	"time"

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
    // Get posts from user database
		userData, err := database.DatabaseService.GetDocument("cyansky-main", "user-data", userId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Internal server error: %v", err),
			})
			return
		}
		var decodedUserData map[string]interface{}
		err = userData.Decode(&decodedUserData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to decode user data: %s", err),
			})
			return
		}

		posts := decodedUserData["posts"].([]interface{})
		sort.Slice(posts, func(i, j int) bool {
			timeI, errI := time.Parse(time.RFC3339, posts[i].(map[string]interface{})["$createdAt"].(string))
			timeJ, errJ := time.Parse(time.RFC3339, posts[j].(map[string]interface{})["$createdAt"].(string))
			if errI != nil || errJ != nil {
				return false
			}
			return timeI.After(timeJ)
		})
    for _, postInterface := range posts {
      post := postInterface.(map[string]interface{})

      if post["author"] == nil {
        post["author"] = make(map[string]interface{})
      }

      authorData := post["author"].(map[string]interface{})
      authorData["username"] = decodedUserData["username"]
    }

		c.JSON(http.StatusOK, gin.H{
			"documents": posts,
		})
	} else {
    // Get posts from global database
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
}
