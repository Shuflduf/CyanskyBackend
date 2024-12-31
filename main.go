package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/models"
	"github.com/atreugo/cors"
	"github.com/joho/godotenv"
	"github.com/savsgio/atreugo/v11"
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
	config := atreugo.Config{
    Addr: "127.0.0.1:8000",
	}
	server := atreugo.New(config)

	cors := cors.New(cors.Config{
		AllowedHeaders:   []string{"Content-Type", "application/json"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})
  server.UseAfter(cors)

	server.Path("GET", "/", func(ctx *atreugo.RequestCtx) error {
		println("Hello World")
		return ctx.HTTPResponse(RetrieveDBDocuments())
	})

	server.Path("POST", "/newacc", func(ctx *atreugo.RequestCtx) error {
		var data map[string]interface{}

		// Parse the JSON body
		if err := json.Unmarshal(ctx.PostBody(), &data); err != nil {
			ctx.Error("Invalid JSON", 400)
			return err
		}

		// Process the data (for demonstration, we just log it)
		log.Printf("Received POST data: %+v\n", data)

		// Send a response
		return ctx.JSONResponse(map[string]string{
			"status":  "success",
			"message": "Data received successfully",
		})
	})

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
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
