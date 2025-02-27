package database

import (
	"fmt"
	"os"

	"github.com/appwrite/sdk-for-go/account"
	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
	"github.com/appwrite/sdk-for-go/query"
)

const ProjectId = "6770875d003cd7a26e8e"

var DatabaseService *databases.Databases
var AccountService *account.Account
var AdminClient client.Client

func SetupDB() {
	AdminClient = appwrite.NewClient(
		appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
		appwrite.WithProject(ProjectId),
		appwrite.WithKey(os.Getenv("ADMIN_SECRET")),
	)

	DatabaseService = appwrite.NewDatabases(AdminClient)
  AccountService = appwrite.NewAccount(AdminClient)
}

func RefreshServices() {
  AdminClient = appwrite.NewClient(
    appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
    appwrite.WithProject(ProjectId),
    appwrite.WithKey(os.Getenv("ADMIN_SECRET")),
  )
  DatabaseService = appwrite.NewDatabases(AdminClient)
  AccountService = appwrite.NewAccount(AdminClient)
}

func GetUserData(userId string) map[string]interface{} {
  RefreshServices()
	documentList, err := DatabaseService.ListDocuments(
		"cyansky-main",
		"user-data",
		DatabaseService.WithListDocumentsQueries([]string{query.Equal("auth-id", userId)}),
	)

	if err != nil {
		fmt.Printf("DB: %v", err)
		return nil
	}

	var info map[string]interface{}
	err = documentList.Decode(&info)
	if err != nil {
		return nil
	}
	return info["documents"].([]interface{})[0].(map[string]interface{})
}
