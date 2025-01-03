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

  database.RefreshServices()

	content := reqBody["content"].(string)
	author := reqBody["author"].(string)
  parent, isReply := reqBody["reply_to"].(string)
  if isReply {
    // get prev replies 
    parentPost, err := database.DatabaseService.GetDocument("cyansky-main", "posts", parent)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error": "Failed to get parent post",
      })
      return
    }
    replyPostId := id.Unique()
    var parentPostData map[string]interface{}
    err = parentPost.Decode(&parentPostData)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error": "Failed to decode parent post",
      })
      return
    }
    parentPostReplies := parentPostData["replies"].([]string)
    parentPostReplies = append(parentPostReplies, replyPostId)
    // edit parent post to add reply 
    _, err = database.DatabaseService.UpdateDocument(
      "cyansky-main",
      "posts",
      parent,
      database.DatabaseService.WithUpdateDocumentData(map[string]interface{}{
        "replies": parentPostReplies,
      }),
    )
  }

	_, err := database.DatabaseService.CreateDocument(
		"cyansky-main",
		"posts",
		id.Unique(),
		map[string]interface{}{
			"content": content,
			"author":  author,
		},
	)

  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": "Failed to create post",
    })
    return
  }

	c.JSON(http.StatusOK, gin.H{
		"message": "Post created",
	})
}
