package main

import (
	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/models"
)

type Document struct {
	Risk string `json:"risk"`
}

type DocumentList struct {
	*models.DocumentList
	Documents []Document `json:"documents"`
}

func main() {
	client := appwrite.NewClient(
		appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
		appwrite.WithProject("6770875d003cd7a26e8e"),
	)

	databases := appwrite.NewDatabases(client)

	response, _ := databases.ListDocuments(
		"67709112001c053f6cdf",
		"677091380032b5cf769d",
	)

	var docs DocumentList
	response.Decode(&docs)

	for _, doc := range docs.Documents {
		println(doc.Risk)
	}
}
