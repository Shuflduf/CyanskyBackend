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
			"error": "Failed to get account data",
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

	// unlike if already liked
	if slices.Contains(likedPosts, postId) && isLike {
		newLikedPosts := slices.DeleteFunc(likedPosts, func(s string) bool {
			return s == postId
		})
		UpdatePostLikes(postId, "likes", -1)
		_, err := database.DatabaseService.UpdateDocument(
			"cyansky-main",
			"user-data",
			authorData["$id"].(string),
			database.DatabaseService.WithUpdateDocumentData(map[string]interface{}{
				"liked-posts": newLikedPosts,
			}),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to update user data: %s", err),
			})
			return
		}

		// undislike if already disliked
	} else if slices.Contains(dislikedPosts, postId) && !isLike {
		newDislikedPosts := slices.DeleteFunc(dislikedPosts, func(s string) bool {
			return s == postId
		})
		UpdatePostLikes(postId, "dislikes", -1)
		_, err := database.DatabaseService.UpdateDocument(
			"cyansky-main",
			"user-data",
			authorData["$id"].(string),
			database.DatabaseService.WithUpdateDocumentData(map[string]interface{}{
				"disliked-posts": newDislikedPosts,
			}),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to update user data: %s", err),
			})
			return
		}

		// like or dislike if not already liked or disliked
	} else if !slices.Contains(likedPosts, postId) && !slices.Contains(dislikedPosts, postId) {
		if isLike {
			likedPosts = append(likedPosts, postId)
			UpdatePostLikes(postId, "likes", 1)
		} else {
			dislikedPosts = append(dislikedPosts, postId)
			UpdatePostLikes(postId, "dislikes", 1)
		}
		_, err := database.DatabaseService.UpdateDocument(
			"cyansky-main",
			"user-data",
			authorData["$id"].(string),
			database.DatabaseService.WithUpdateDocumentData(map[string]interface{}{
				"liked-posts":    likedPosts,
				"disliked-posts": dislikedPosts,
			}),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to update user data: %s", err),
			})
			return
		}

		// like if already disliked
	}

	c.JSON(http.StatusOK, gin.H{
		"liked_posts": likedPosts,
	})
}

func UpdateUserLiked(toChange, field, userId string) {
	user, err := database.DatabaseService.GetDocument("cyansky-main", "user-data", userId)
	if err != nil {
		fmt.Println(err)
		return
	}
	var userData map[string]interface{}
	err = user.Decode(&userData)
	if err != nil {
		fmt.Println(err)
		return
	}
	likedPosts := userData[field].([]interface{})
	dislikedPosts := userData["disliked-posts"].([]interface{})

	// if toChange is in field, remove it 

	// if toChange is in the opposite field, remove it and add it to field 

	// if toChange is not in either field, add it to field

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
		return
	}
}

func UpdatePostLikes(postId string, field string, change int) {
	database.RefreshServices()
	post, err := database.DatabaseService.GetDocument("cyansky-main", "posts", postId)
	if err != nil {
		fmt.Println(err)
		return
	}
	var postData map[string]interface{}
	err = post.Decode(&postData)
	if err != nil {
		fmt.Println(err)
		return
	}
	likes := postData[field].(float64)
	likes += float64(change)
	postData[field] = likes
	_, err = database.DatabaseService.UpdateDocument(
		"cyansky-main",
		"posts",
		postId,
		database.DatabaseService.WithUpdateDocumentData(map[string]interface{}{
			field: likes,
		}),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func convertToStringSlice(input []interface{}) []string {
	output := make([]string, len(input))
	for i, v := range input {
		output[i] = v.(string)
	}
	return output
}
