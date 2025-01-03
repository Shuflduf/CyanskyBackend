package routes

import (
	"cyansky/appwrite"
	"fmt"
	"net/http"

	"github.com/appwrite/sdk-for-go/appwrite"
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
	// author := reqBody["author"].(string)
  authorClient := appwrite.NewClient(
    appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
    appwrite.WithProject(database.ProjectId),
    appwrite.WithSession(reqBody["token"].(string)),
  )
  service := appwrite.NewAccount(authorClient)
  accountData, accErr := service.Get()
  if accErr != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
      "error": "Failed to get account data",
    })
    return
  }
  authorAuthId := accountData.Id
  authorData := database.GetUserData(authorAuthId)
  author := authorData["$id"].(string)

  newPostId := id.Unique()
  parent, isReply := reqBody["reply_to"].(string)
  if isReply {
    parentPost, err := database.DatabaseService.GetDocument("cyansky-main", "posts", parent)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error": fmt.Sprintf("Failed to get parent post: %s", err),
      })
      return
    }
    var parentPostData map[string]interface{}
    err = parentPost.Decode(&parentPostData)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error": "Failed to decode parent post",
      })
      return
    }
    parentPostReplies := parentPostData["replies"].([]interface{})
    parentPostReplies = append(parentPostReplies, newPostId)
    // edit parent post to add reply 
    _, err = database.DatabaseService.UpdateDocument(
      "cyansky-main",
      "posts",
      parent,
      database.DatabaseService.WithUpdateDocumentData(map[string]interface{}{
        "replies": parentPostReplies,
      }),
    )
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error": fmt.Sprintf("Failed to update parent post: %s", err),
      })
      return
    }
  }

	_, err := database.DatabaseService.CreateDocument(
		"cyansky-main",
		"posts",
		newPostId,
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
