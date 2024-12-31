package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Document struct {
	Risk string `json:"risk"`
}

type DocumentList struct {
	*models.DocumentList
	Documents []Document `json:"documents"`
}

var myDatabases *databases.Databases

func main() {
	_ = godotenv.Load()
	SetupDB()
	SetupServer()
}

func SetupServer() {
	r := gin.Default()

  r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"*"},
    AllowMethods:     []string{"GET", "POST"},
    AllowHeaders:     []string{"Origin", "Content-Type"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    AllowOriginFunc: func(origin string) bool {
      return origin == "https://github.com"
    },
    MaxAge: 12 * time.Hour,
  }))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": RetrieveDBDocuments(),
		})
	})
	r.POST("/newacc", func(c *gin.Context) {
		var reqBody map[string]any
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		fmt.Println(reqBody)

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": reqBody,
		})
	})

	r.Run(":8000")
}

func SetupDB() {
	client := appwrite.NewClient(
		appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
		appwrite.WithProject("6770875d003cd7a26e8e"),
		appwrite.WithKey(os.Getenv("ADMIN_SECRET")),
	)

	myDatabases = appwrite.NewDatabases(client)
}

func RetrieveDBDocuments() string {
	response, _ := myDatabases.ListDocuments(
		"67709112001c053f6cdf",
		"677091380032b5cf769d",
	)

	var docs DocumentList
	response.Decode(&docs)

	return docs.Documents[0].Risk
}
