package routes

import (
	"cyansky/appwrite"
	"fmt"
	"net/http"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/gin-gonic/gin"
)

func Follow(c *gin.Context) {
	var reqBody map[string]any
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}
	database.RefreshServices()

	userClient := appwrite.NewClient(
		appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
		appwrite.WithProject(database.ProjectId),
		appwrite.WithSession(reqBody["token"].(string)),
	)

	userAccount := appwrite.NewAccount(userClient)
	user, err := userAccount.Get()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to get user: %v", err),
		})
		return
	}
	userData := database.GetUserData(user.Id)
	userDataId := userData["$id"].(string)

  // get current followers
  document, err := database.DatabaseService.GetDocument(
    "cyansky-main",
    "user-data",
    reqBody["target"].(string),
  )
  var targetUserData map[string]interface{}
  err = document.Decode(&targetUserData)
  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error": fmt.Sprintf("Failed to get user data: %v", err),
    })
    return
  }
  targetCurrentFollowers := targetUserData["following"].([]interface{})
  followed := false
  for i, follower := range targetCurrentFollowers {
    if follower == userDataId {
      followed = true
      targetCurrentFollowers = append(targetCurrentFollowers[:i], targetCurrentFollowers[i+1:]...)
      break
    }
  }
  if !followed {
    targetCurrentFollowers = append(targetCurrentFollowers, userDataId)
  }

	database.RefreshServices()
	_, err = database.DatabaseService.UpdateDocument(
		"cyansky-main",
		"user-data",
		reqBody["target"].(string),
		database.DatabaseService.WithUpdateDocumentData(map[string]interface{}{
      "following": targetCurrentFollowers,
    }),
	)

  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{
      "error": fmt.Sprintf("Failed to update user data: %v", err),
    })
    return
  }

  if followed {
    c.JSON(http.StatusOK, gin.H{
      "message": "Unfollowed successfully",
    })
  } else {
    c.JSON(http.StatusOK, gin.H{
      "message": "Followed successfully",
    })
  }
}
