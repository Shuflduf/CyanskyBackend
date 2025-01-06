package routes

import (
	"net/http"

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

	// isLike := reqBody["is_like"].(bool)
}
