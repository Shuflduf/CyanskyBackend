package routes

import (
	"cyansky/appwrite"
	"net/http"

	"github.com/appwrite/sdk-for-go/id"
	"github.com/gin-gonic/gin"
)

func MakePost(c *gin.Context) {
	var reqBody map[string]any
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	content := reqBody["content"].(string)
	author := reqBody["author"].(string)

	_, err := database.DatabaseManager.CreateDocument(
		"cyansky-main",
		"posts",
		id.Unique(),
		map[string]interface{}{
			"content": content,
			"author":  author,
		},
	)

  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error": "Failed to create post",
    })
    return
  }

	c.JSON(http.StatusOK, gin.H{
		"message": "Post created",
	})
}
