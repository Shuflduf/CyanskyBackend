package routes

import (
	"cyansky/appwrite"
	"fmt"
	"net/http"
	"slices"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/gin-gonic/gin"
)

func Likes(c *gin.Context) {
	var reqBody map[string]any
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	database.RefreshServices()

	isLike := reqBody["is_like"].(bool)
	var likeValue int
	if isLike {
		likeValue = 1
	} else {
		likeValue = -1
	}
	postId := reqBody["post_id"].(string)
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
	authorAuthId := accountData.Id
	authorData := database.GetUserData(authorAuthId)

	likedPostsInterface := authorData["liked-posts"].([]interface{})
	dislikedPostsInterface := authorData["disliked-posts"].([]interface{})
	likedPosts := convertToStringSlice(likedPostsInterface)
	dislikedPosts := convertToStringSlice(dislikedPostsInterface)

	if slices.Contains(likedPosts, postId) && slices.Contains(dislikedPosts, postId) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Post is both liked and disliked",
		})
		return
	}

	result, err := UpdateUserLiked(postId, authorData["$id"].(string), likeValue)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to update user data: %s", err),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func UpdateUserLiked(postId, userId string, likeValue int) (map[string]interface{}, error) {
	user, err := database.DatabaseService.GetDocument("cyansky-main", "user-data", userId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var userData map[string]interface{}
	err = user.Decode(&userData)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	likedPosts := convertToStringSlice(userData["liked-posts"].([]interface{}))
	dislikedPosts := convertToStringSlice(userData["disliked-posts"].([]interface{}))

	// coudlve been a function
	removePost := func(slice *[]string, postId string) {
		if slices.Contains(*slice, postId) {
			*slice = slices.DeleteFunc(*slice, func(s string) bool {
				return s == postId
			})
		}
	}

  // if i want to LIKE
	if likeValue == 1 {
    // remove DISLIKE
    if slices.Contains(dislikedPosts, postId) {
      removePost(&dislikedPosts, postId)
      UpdatePostLikes(postId, "dislikes", -1)
    }
    // add LIKE
    if !slices.Contains(likedPosts, postId) {
      likedPosts = append(likedPosts, postId)
      UpdatePostLikes(postId, "likes", 1)
    } else {
      // if already LIKED, remove LIKE
      removePost(&likedPosts, postId)
      UpdatePostLikes(postId, "likes", -1)
    }
  } else if likeValue == -1 {
    // remove LIKE 
    if slices.Contains(likedPosts, postId) {
      removePost(&likedPosts, postId)
      UpdatePostLikes(postId, "likes", -1)
    }
    // add DISLIKE 
    if !slices.Contains(dislikedPosts, postId) {
      dislikedPosts = append(dislikedPosts, postId)
      UpdatePostLikes(postId, "dislikes", 1)
    } else {
      // if already DISLIKED, remove DISLIKE 
      removePost(&dislikedPosts, postId)
      UpdatePostLikes(postId, "dislikes", -1)
    }
  }

	_, err = database.DatabaseService.UpdateDocument(
		"cyansky-main",
		"user-data",
		userId,
		database.DatabaseService.WithUpdateDocumentData(map[string]interface{}{
			"liked-posts":    likedPosts,
			"disliked-posts": dislikedPosts,
		}),
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return map[string]interface{}{
		"liked_posts":    likedPosts,
		"disliked_posts": dislikedPosts,
		"new_likes":      GetLikes(false, postId),
		"new_dislikes":   GetLikes(true, postId),
	}, nil
}

func GetLikes(isDislike bool, postId string) int {
	database.RefreshServices()
	post, err := database.DatabaseService.GetDocument("cyansky-main", "posts", postId)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	var postData map[string]interface{}
	err = post.Decode(&postData)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if isDislike {
		return int(postData["dislikes"].(float64))
	} else {
		return int(postData["likes"].(float64))
	}
}

func UpdatePostLikes(postId string, field string, change int) int {
	database.RefreshServices()

	likes := float64(GetLikes(field == "dislikes", postId))
	likes += float64(change)
	// postData[field] = likes
	_, err := database.DatabaseService.UpdateDocument(
		"cyansky-main",
		"posts",
		postId,
		database.DatabaseService.WithUpdateDocumentData(map[string]interface{}{
			field: likes,
		}),
	)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return int(likes)
}

func convertToStringSlice(input []interface{}) []string {
	output := make([]string, len(input))
	for i, v := range input {
		output[i] = v.(string)
	}
	return output
}
