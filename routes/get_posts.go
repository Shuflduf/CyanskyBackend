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

	var postList map[string]interface{}
	var decodedDocuments map[string]interface{}
	err = documentList.Decode(&decodedDocuments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to decode documents: %s", err),
		})
	}

  decodedDocumentList := decodedDocuments["documents"].([]interface{})
  for _, document := range decodedDocumentList {
    database.RefreshServices()
    documentId := document.(map[string]interface{})["$id"].(string)
    DocumentData, err := database.DatabaseService.GetDocument(
      "cyansky-main",
      "posts",
      documentId,
      )
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error": fmt.Sprintf("Internal server error: %v", err),
      })
      return
    }
    var decodedDocumentData map[string]interface{}
    err = DocumentData.Decode(&decodedDocumentData)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{
        "error": fmt.Sprintf("Failed to decode document data: %s", err),
      })
    } 
    if postList == nil {
      postList = make(map[string]interface{})
      postList["posts"] = []interface{}{}
    }
    postList["posts"] = append(postList["posts"].([]interface{}), decodedDocumentData)

  }

	c.JSON(http.StatusOK, postList)
}
