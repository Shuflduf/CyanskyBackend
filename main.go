package main

import (
	"cyansky/appwrite"
	"cyansky/routes"
	"net/http"
	"time"

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


func main() {
	_ = godotenv.Load()
	database.SetupDB()
	SetupServer()
}

func SetupServer() {
	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"192.168.1.0/24"})

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "🐸🚀",
		})
	})

	r.POST("/createpost", routes.MakePost)
  r.POST("/createaccount", routes.CreateAccount)

	r.Run(":8000")
}


// func RetrieveDBDocuments() string {
// 	response, _ := databaseManager.ListDocuments(
// 		"67709112001c053f6cdf",
// 		"677091380032b5cf769d",
// 	)
//
// 	var docs DocumentList
// 	response.Decode(&docs)
//
// 	return docs.Documents[0].Risk
// }
