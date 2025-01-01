package database

import (
	"os"

	"github.com/appwrite/sdk-for-go/account"
	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/databases"
)

var DatabaseManager *databases.Databases
var AccountManager *account.Account
var BaseClient client.Client

func SetupDB() {
	BaseClient = appwrite.NewClient(
		appwrite.WithEndpoint("https://cloud.appwrite.io/v1"),
		appwrite.WithProject("6770875d003cd7a26e8e"),
		appwrite.WithKey(os.Getenv("ADMIN_SECRET")),
	)

	DatabaseManager = appwrite.NewDatabases(BaseClient)
}
