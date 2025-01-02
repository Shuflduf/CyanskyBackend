package database

import (
	"os"

	"github.com/appwrite/sdk-for-go/account"
	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
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

// func GetAccountInfo(authId string) interface{} {
// 	searchQuery := []string{query.Equal("auth-id", []interface{}{authId})}
// 	documentList, err := DatabaseService.ListDocuments(
// 		"cyansky-main",
// 		"user-data",
// 		DatabaseService.WithListDocumentsQueries(searchQuery),
// 	)
//
//   if err != nil {
//     fmt.Println(err)
//     return nil
//   }
//
//   var info interface{}
//   return documentList.Documents[0].Decode(info)
// }
