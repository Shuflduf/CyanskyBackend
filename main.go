package main

import (
	"cyansky/appwrite"
	"cyansky/routes"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	database.SetupDB()
	SetupServer()
}

func SetupServer() {
	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"192.168.1.0/24"})

  // i hate cors
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
  r.POST("/login", routes.Login)
  r.POST("/follow", routes.Follow)

	r.Run(":8000")
}
